-- ## Indexes
DROP INDEX IF EXISTS idx_class_users;

DROP INDEX IF EXISTS idx_submissions;

DROP INDEX IF EXISTS idx_grades;

DROP INDEX IF EXISTS idx_grade_criteria;

DROP INDEX IF EXISTS idx_users_fullname;

-- ## Tables
DROP TABLE IF EXISTS grade_has_criteria;

DROP TABLE IF EXISTS grades;

DROP TABLE IF EXISTS submissions;

DROP TABLE IF EXISTS test_blocks;

DROP TABLE IF EXISTS languages;

DROP TABLE IF EXISTS markdown_blocks;

DROP TABLE IF EXISTS laboratories;

DROP TABLE IF EXISTS criteria;

DROP TABLE IF EXISTS objectives;

DROP TABLE IF EXISTS rubrics;

DROP TABLE IF EXISTS class_has_users;

DROP TABLE IF EXISTS classes;

DROP TABLE IF EXISTS colors;

DROP TABLE IF EXISTS users;

-- ## Types
DROP TYPE IF EXISTS SUBMISSION_STATUS;

DROP TYPE IF EXISTS USER_ROLES;

-- ## Extensions
DROP EXTENSION IF EXISTS "uuid-ossp";