-- ## Triggers
DROP TRIGGER IF EXISTS set_created_by ON users;
DROP FUNCTION IF EXISTS update_created_by();

-- ## Indexes
DROP INDEX IF EXISTS idx_class_users;

DROP INDEX IF EXISTS idx_submissions;

DROP INDEX IF EXISTS idx_grades;

DROP INDEX IF EXISTS idx_grade_criteria;

DROP INDEX IF EXISTS idx_blocks_index;

DROP INDEX IF EXISTS idx_users_lower_fullName;

DROP INDEX IF EXISTS idx_users_role;

-- ## Views
DROP VIEW IF EXISTS submissions_work_metadata;

DROP VIEW IF EXISTS courses_with_color;

DROP VIEW IF EXISTS courses_has_users_view;

DROP VIEW IF EXISTS users_with_creator;

DROP VIEW IF EXISTS objectives_owners;

DROP VIEW IF EXISTS criteria_owners;

-- ## Tables
DROP TABLE IF EXISTS archives;

DROP TABLE IF EXISTS grade_has_criteria;

DROP TABLE IF EXISTS grades;

DROP TABLE IF EXISTS submissions;

DROP TABLE IF EXISTS blocks_index;

DROP TABLE IF EXISTS test_blocks;

DROP TABLE IF EXISTS languages;

DROP TABLE IF EXISTS markdown_blocks;

DROP TABLE IF EXISTS laboratories;

DROP TABLE IF EXISTS criteria;

DROP TABLE IF EXISTS objectives;

DROP TABLE IF EXISTS rubrics;

DROP TABLE IF EXISTS courses_has_users;

DROP TABLE IF EXISTS invitation_codes;

DROP TABLE IF EXISTS courses;

DROP TABLE IF EXISTS colors;

DROP TABLE IF EXISTS users;

-- ## Types
DROP TYPE IF EXISTS SUBMISSION_STATUS;

DROP TYPE IF EXISTS USER_ROLES;

-- ## Extensions
DROP EXTENSION IF EXISTS "uuid-ossp";
DROP EXTENSION IF EXISTS "citext";