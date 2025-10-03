-- +migrate Down

DROP TABLE IF EXISTS applications;
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS job_seeker_profiles;
DROP TABLE IF EXISTS employer_profiles;
DROP TABLE IF EXISTS users;

