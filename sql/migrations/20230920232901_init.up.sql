-- ## Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE EXTENSION IF NOT EXISTS citext;

-- ## Types
CREATE TYPE SUBMISSION_STATUS AS ENUM ('pending', 'running', 'ready');

CREATE TYPE USER_ROLES AS ENUM ('admin', 'teacher', 'student');

-- ## Tables 
CREATE TABLE IF NOT EXISTS users (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "role" USER_ROLES NOT NULL DEFAULT 'student',
  "institutional_id" VARCHAR(16) NULL UNIQUE,
  "email" CITEXT NOT NULL UNIQUE,
  "full_name" VARCHAR NOT NULL,
  "password_hash" VARCHAR NOT NULL,
  "created_by" UUID NULL REFERENCES users(id),
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS colors(
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "hexadecimal" VARCHAR(9) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS courses (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "teacher_id" UUID NOT NULL REFERENCES users(id),
  "color_id" UUID NOT NULL REFERENCES colors(id),
  "name" VARCHAR(96) NOT NULL
);

CREATE TABLE IF NOT EXISTS invitation_codes (
  "course_id" UUID PRIMARY KEY REFERENCES courses(id),
  "code" VARCHAR(9) NOT NULL UNIQUE CHECK (LENGTH(code) >= 9),
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS courses_has_users (
  "course_id" UUID NOT NULL REFERENCES courses(id),
  "user_id" UUID NOT NULL REFERENCES users(id),
  "is_class_hidden" BOOLEAN NOT NULL DEFAULT FALSE,
  "is_user_active" BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS rubrics (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "teacher_id" UUID NOT NULL REFERENCES users(id),
  "name" VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS objectives (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "rubric_id" UUID NOT NULL REFERENCES rubrics(id),
  "description" VARCHAR(510) NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS criteria (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "objective_id" UUID NOT NULL REFERENCES objectives(id) ON DELETE CASCADE,
  "description" VARCHAR(510) NOT NULL,
  "weight" DECIMAL(9, 6) NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS laboratories (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "course_id" UUID NOT NULL REFERENCES courses(id),
  "rubric_id" UUID DEFAULT NULL REFERENCES rubrics(id) ON DELETE
  SET
    DEFAULT,
    "name" VARCHAR(255) NOT NULL,
    "opening_date" TIMESTAMP NOT NULL,
    "due_date" TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS blocks_index (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "laboratory_id" UUID NOT NULL REFERENCES laboratories(id),
  "block_position" SMALLINT NOT NULL
);

CREATE TABLE IF NOT EXISTS markdown_blocks (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "laboratory_id" UUID NOT NULL REFERENCES laboratories(id),
  "block_index_id" UUID NOT NULL REFERENCES blocks_index(id) ON DELETE CASCADE,
  "content" TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS archives (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "file_id" UUID NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS languages (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "template_archive_id" UUID NOT NULL UNIQUE REFERENCES archives(id),
  "name" VARCHAR(32) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS test_blocks (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "language_id" UUID NOT NULL REFERENCES languages(id),
  "test_archive_id" UUID NOT NULL UNIQUE REFERENCES archives(id),
  "laboratory_id" UUID NOT NULL REFERENCES laboratories(id),
  "block_index_id" UUID NOT NULL REFERENCES blocks_index(id) ON DELETE CASCADE,
  "name" VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS submissions (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "test_id" UUID NOT NULL REFERENCES test_blocks(id),
  "student_id" UUID NOT NULL REFERENCES users(id),
  "archive_id" UUID NOT NULL UNIQUE REFERENCES archives(id),
  "passing" BOOLEAN NOT NULL DEFAULT FALSE,
  "status" SUBMISSION_STATUS NOT NULL DEFAULT 'pending',
  "stdout" TEXT NOT NULL DEFAULT '',
  "submitted_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS grades (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "laboratory_id" UUID NOT NULL REFERENCES laboratories(id),
  "student_id" UUID NOT NULL REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS grade_has_criteria (
  "grade_id" UUID NOT NULL REFERENCES grades(id),
  "criteria_id" UUID NOT NULL REFERENCES criteria(id) ON DELETE CASCADE,
  "objective_id" UUID NOT NULL REFERENCES objectives(id) ON DELETE CASCADE
);

-- ## Indexes
-- ### Unique indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_class_users ON courses_has_users(course_id, user_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_submissions ON submissions(test_id, student_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_grades ON grades(laboratory_id, student_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_grade_criteria ON grade_has_criteria(grade_id, objective_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_blocks_index ON blocks_index(laboratory_id, block_position);

-- ### Search indexes
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

CREATE INDEX IF NOT EXISTS idx_users_lower_fullName ON users(LOWER(full_name));

-- ## Views
--- ### Users
CREATE
OR REPLACE VIEW users_with_creator AS
SELECT
  users.id,
  users.role,
  users.institutional_id,
  users.email,
  users.full_name,
  users.created_by,
  creator.full_name AS creator_full_name,
  users.created_at
FROM
  users
  LEFT JOIN users AS creator ON users.created_by = creator.id;

--- ### courses
CREATE
OR REPLACE VIEW courses_with_color AS
SELECT
  courses.id,
  courses.teacher_id,
  courses.name,
  colors.hexadecimal AS color
FROM
  courses
  INNER JOIN colors ON courses.color_id = colors.id;

--- ### courses_has_users
CREATE
OR REPLACE VIEW courses_has_users_view AS
SELECT
  courses_has_users.course_id,
  courses.name AS course_name,
  courses.teacher_id AS course_teacher_id,
  colors.hexadecimal AS course_color,
  courses_has_users.user_id,
  users.full_name AS user_full_name,
  users.email AS user_email,
  users.role AS user_role,
  users.institutional_id AS user_institutional_id,
  courses_has_users.is_class_hidden,
  courses_has_users.is_user_active
FROM
  courses_has_users
  INNER JOIN users ON courses_has_users.user_id = users.id
  INNER JOIN courses ON courses_has_users.course_id = courses.id
  INNER JOIN colors ON courses.color_id = colors.id;

--- ### Objectives
CREATE
OR REPLACE VIEW objectives_owners AS
SELECT
  objectives.id AS objective_id,
  rubrics.teacher_id
FROM
  objectives
  INNER JOIN rubrics ON objectives.rubric_id = rubrics.id;

--- ### Criteria
CREATE
OR REPLACE VIEW criteria_owners AS
SELECT
  criteria.id AS criteria_id,
  rubrics.teacher_id
FROM
  criteria
  INNER JOIN objectives ON criteria.objective_id = objectives.id
  INNER JOIN rubrics ON objectives.rubric_id = rubrics.id;

-- ## Triggers
--- ### Update created_by on users
CREATE
OR REPLACE FUNCTION update_created_by()
RETURNS TRIGGER 
LANGUAGE PLPGSQL
AS $$
BEGIN
  IF NEW.created_by IS NULL THEN
    NEW.created_by := NEW.id;
  END IF;

  RETURN NEW;
END $$
;

CREATE
OR REPLACE TRIGGER set_created_by BEFORE
INSERT
  ON users FOR EACH ROW EXECUTE PROCEDURE update_created_by();

-- ## Data
-- ### Colors
INSERT INTO
  colors (hexadecimal)
VALUES
  ('#a78bfa'),
  ('#34d399'),
  ('#f87171'),
  ('#22d3ee'),
  ('#fbbf24'),
  ('#f472b6');

-- ### Languages
DO $$

DECLARE 
  JAVA_FILESYSTEM_ARCHIVE_UUID UUID;
  JAVA_DB_ARCHIVE_UUID UUID;

BEGIN 
  JAVA_FILESYSTEM_ARCHIVE_UUID := '487034c9-441c-4fb9-b0f3-8f4dd6176532';

  INSERT INTO
    archives (file_id)
  VALUES
    (JAVA_FILESYSTEM_ARCHIVE_UUID)
  RETURNING
    id
  INTO
    JAVA_DB_ARCHIVE_UUID;

  INSERT INTO
    languages (name, template_archive_id)
  VALUES
    (
      'Java',
      JAVA_DB_ARCHIVE_UUID
    );

END $$;

-- ### Admin user (To be used in development)
INSERT INTO
  users (
    role,
    email,
    full_name,
    password_hash
  )
VALUES
  (
    'admin',
    'development.admin@gmail.com',
    'Development Admin',
    '$argon2id$v=19$m=16,t=1,p=1$UUVzSDRZQkpKZkhrWmN4ZA$TiQHkBQI7A+1987WmMDHhw'
  );