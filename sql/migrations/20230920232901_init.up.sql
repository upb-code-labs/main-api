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
  "created_by" UUID DEFAULT NULL REFERENCES users(id),
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
  "course_id" UUID PRIMARY KEY REFERENCES courses(id) ON DELETE CASCADE,
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
  "rubric_id" UUID NOT NULL REFERENCES rubrics(id) ON DELETE CASCADE,
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
  "rubric_id" UUID DEFAULT NULL REFERENCES rubrics(id) ON DELETE SET DEFAULT,
  "name" VARCHAR(255) NOT NULL,
  "opening_date" TIMESTAMP WITH TIME ZONE NOT NULL,
  "due_date" TIMESTAMP WITH TIME ZONE NOT NULL
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
  "test_block_id" UUID NOT NULL REFERENCES test_blocks(id) ON DELETE CASCADE,
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
  "rubric_id" UUID NOT NULL REFERENCES rubrics(id) ON DELETE CASCADE,
  "student_id" UUID NOT NULL REFERENCES users(id),
  "comment" TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS grade_has_criteria (
  "grade_id" UUID NOT NULL REFERENCES grades(id),
  "criteria_id" UUID NULL NULL REFERENCES criteria(id) ON DELETE SET NULL,
  "objective_id" UUID NOT NULL REFERENCES objectives(id) ON DELETE CASCADE
);

-- ## Indexes
-- ### Unique indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_class_users ON courses_has_users(course_id, user_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_submissions ON submissions(test_block_id, student_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_grades ON grades(laboratory_id, rubric_id, student_id);

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

-- ### Submissions work to be sent to the work queue
CREATE OR REPLACE VIEW submissions_work_metadata AS
SELECT
  submissions.id AS submission_id,
  language_archive.file_id AS language_file_id,
  test_archive.file_id AS test_file_id,
  submission_archive.file_id AS submission_file_id
FROM submissions 
  INNER JOIN test_blocks ON submissions.test_block_id = test_blocks.id
  INNER JOIN languages ON test_blocks.language_id = languages.id
  INNER JOIN archives AS language_archive ON languages.template_archive_id = language_archive.id
  INNER JOIN archives AS test_archive ON test_blocks.test_archive_id = test_archive.id
  INNER JOIN archives AS submission_archive ON submissions.archive_id = submission_archive.id;

-- ### Students progress
CREATE OR REPLACE VIEW students_progress_view AS
SELECT
  users.id AS student_id,
  users.full_name as student_full_name,
  test_blocks.laboratory_id,
  COUNT(submissions.id) FILTER (
	  WHERE submissions.status = 'pending'
  ) AS pending_submissions,
  COUNT(submissions.id) FILTER (
	  WHERE submissions.status = 'running'
  ) AS running_submissions,
  COUNT(submissions.id) FILTER (
	  WHERE submissions.status = 'ready' AND submissions.passing = FALSE
  ) AS failing_submissions,
  COUNT(submissions.id) FILTER (
	  WHERE submissions.status = 'ready' AND submissions.passing = TRUE
  ) AS success_submissions
FROM
  submissions
JOIN
  users ON submissions.student_id = users.id
JOIN
  test_blocks ON submissions.test_block_id = test_blocks.id
GROUP BY
  users.id, users.full_name, test_blocks.laboratory_id;


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

-- ### Summarized grades
CREATE
OR REPLACE VIEW summarized_grades AS
SELECT
  grades.id AS grade_id,
  grades.student_id,
  students.full_name AS student_full_name,
  grades.laboratory_id,
  grades.rubric_id,
  SUM(criteria.weight) AS total_criteria_weight,
  grades.comment
FROM
  grades
  INNER JOIN users AS students ON grades.student_id = students.id
  INNER JOIN grade_has_criteria ON grades.id = grade_has_criteria.grade_id
  INNER JOIN criteria ON grade_has_criteria.criteria_id = criteria.id
GROUP BY
  grades.id, students.full_name;

-- ## Procedures and functions
--- ### Swap blocks index
CREATE
OR REPLACE FUNCTION swap_blocks_index(
  IN first_block_id UUID,
  IN second_block_id UUID
)
RETURNS VOID
LANGUAGE PLPGSQL
AS $$
DECLARE
  is_first_block_a_markdown_block BOOLEAN;
  is_second_block_a_markdown_block BOOLEAN;
  first_block_index_id UUID;
  second_block_index_id UUID;
BEGIN
  -- Get the type of the blocks
  SELECT
    EXISTS(
      SELECT
        1
      FROM
        markdown_blocks
      WHERE
        id = first_block_id
    )
  INTO
    is_first_block_a_markdown_block;

  SELECT
    EXISTS(
      SELECT
        1
      FROM
        markdown_blocks
      WHERE
        id = second_block_id
    )
  INTO
    is_second_block_a_markdown_block;
  
  -- Get the index of the blocks
  IF is_first_block_a_markdown_block THEN
    SELECT
      block_index_id
    INTO
      first_block_index_id
    FROM
      markdown_blocks
    WHERE
      id = first_block_id;
  ELSE
    SELECT
      block_index_id
    INTO
      first_block_index_id
    FROM
      test_blocks
    WHERE
      id = first_block_id;
  END IF;

  IF is_second_block_a_markdown_block THEN
    SELECT
      block_index_id
    INTO
      second_block_index_id
    FROM
      markdown_blocks
    WHERE
      id = second_block_id;
  ELSE
    SELECT
      block_index_id
    INTO
      second_block_index_id
    FROM
      test_blocks
    WHERE
      id = second_block_id;
  END IF;

  -- Swap the indexes
  IF is_first_block_a_markdown_block THEN
    UPDATE
      markdown_blocks
    SET
      block_index_id = second_block_index_id
    WHERE
      id = first_block_id;
  ELSE
    UPDATE
      test_blocks
    SET
      block_index_id = second_block_index_id
    WHERE
      id = first_block_id;
  END IF;

  IF is_second_block_a_markdown_block THEN
    UPDATE
      markdown_blocks
    SET
      block_index_id = first_block_index_id
    WHERE
      id = second_block_id;
  ELSE
    UPDATE
      test_blocks
    SET
      block_index_id = first_block_index_id
    WHERE
      id = second_block_id;
  END IF;
END $$
;

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
  ('#8b5cf6'),
  ('#ef4444'),
  ('#3b82f6'),
  ('#ea580c'),
  ('#ec4899'),
  ('#16a34a'),
  ('#f43f5e'),
  ('#d946ef'),
  ('#a855f7'),
  ('#6366f1');

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
      'Java JDK 17',
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