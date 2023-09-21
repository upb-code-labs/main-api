-- ## Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ## Tables 
CREATE TABLE IF NOT EXISTS roles (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "name" VARCHAR(32) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "role_id" UUID NOT NULL REFERENCES roles(id),
  "institutional_id" VARCHAR(16) NOT NULL UNIQUE,
  "email" VARCHAR(64) NOT NULL UNIQUE,
  "full_name" VARCHAR NOT NULL,
  "password_hash" VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS colors(
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "hexadecimal" VARCHAR(7) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS classes (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "teacher_id" UUID NOT NULL REFERENCES users(id),
  "color_id" UUID NOT NULL REFERENCES colors(id),
  "invitation_code" VARCHAR(8) NOT NULL UNIQUE,
  "name" VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS class_has_users (
  "class_id" UUID NOT NULL REFERENCES classes(id),
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
  "name" VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS criteria (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "objective_id" UUID NOT NULL REFERENCES objectives(id),
  "description" VARCHAR NOT NULL,
  "weight" DECIMAL(5, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS laboratories (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "class_id" UUID NOT NULL REFERENCES classes(id),
  "rubric_id" UUID NOT NULL REFERENCES rubrics(id),
  "name" VARCHAR(255) NOT NULL,
  "opening_date" TIMESTAMP NOT NULL,
  "due_date" TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS markdown_blocks (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "laboratory_id" UUID NOT NULL REFERENCES laboratories(id),
  "content" TEXT NOT NULL DEFAULT '',
  "order" SMALLINT NOT NULL
);

CREATE TABLE IF NOT EXISTS languages (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "name" VARCHAR(32) NOT NULL UNIQUE,
  "base_archive" BYTEA NOT NULL
);

CREATE TABLE IF NOT EXISTS test_blocks (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "laboratory_id" UUID NOT NULL REFERENCES laboratories(id),
  "language_id" UUID NOT NULL REFERENCES languages(id),
  "name" VARCHAR(255) NOT NULL,
  "tests_archive" BYTEA NOT NULL,
  "order" SMALLINT NOT NULL
);

CREATE TYPE SUBMISSION_STATUS AS ENUM ('pending', 'running', 'ready');

CREATE TABLE IF NOT EXISTS submissions (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "test_id" UUID NOT NULL REFERENCES test_blocks(id),
  "student_id" UUID NOT NULL REFERENCES users(id),
  "archive" BYTEA NOT NULL,
  "passing" BOOLEAN NOT NULL DEFAULT FALSE,
  "status" SUBMISSION_STATUS NOT NULL DEFAULT 'pending',
  "stdout" TEXT NOT NULL DEFAULT '',
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS grades (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "laboratory_id" UUID NOT NULL REFERENCES laboratories(id),
  "student_id" UUID NOT NULL REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS grade_has_criteria (
  "grade_id" UUID NOT NULL REFERENCES grades(id),
  "criteria_id" UUID NOT NULL REFERENCES criteria(id),
  "objective_id" UUID NOT NULL REFERENCES objectives(id)
);

-- ## Indexes
-- ### Unique indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_class_users ON class_has_users(class_id, user_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_submissions ON submissions(test_id, student_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_grades ON grades(laboratory_id, student_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_grade_criteria ON grade_has_criteria(grade_id, objective_id);

-- ### Search indexes
CREATE INDEX IF NOT EXISTS idx_users_fullname ON users(full_name);

-- ## Data
-- ### Roles
INSERT INTO
  roles (name)
VALUES
  ('admin'),
  ('teacher'),
  ('student');