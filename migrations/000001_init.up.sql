-- +migrate Up

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('job_seeker', 'employer', 'admin')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS employer_profiles (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    company_name VARCHAR(255) NOT NULL,
    industry VARCHAR(100),
    website VARCHAR(255),
    description TEXT,
    logo_url VARCHAR(255),
    location VARCHAR(255),
    company_size VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS job_seeker_profiles (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    headline VARCHAR(255),
    summary TEXT,
    phone VARCHAR(20),
    location VARCHAR(255),
    resume_url VARCHAR(255),
    logo_url VARCHAR(255),
    skills TEXT[],
    experience_level VARCHAR(50)  DEFAULT 'entry' CHECK (experience_level IN ('Entry-level', 'Mid-level','Senior', 'Lead')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    employer_id INT NOT NULL REFERENCES employer_profiles(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    location VARCHAR(255),
    job_type VARCHAR(50) CHECK (job_type IN ('full_time', 'part_time', 'contract', 'internship', 'remote')),
    salary_min DECIMAL(12, 2) DEFAULT 1000.00,
    salary_max DECIMAL(12, 2) DEFAULT 1000.00,
    experience_level VARCHAR(50)  DEFAULT 'entry' CHECK (experience_level IN ('Entry-level', 'Mid-level','Senior', 'Lead')),
    required_skills TEXT[],
    category VARCHAR(100),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'closed', 'draft')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS applications (
    id SERIAL PRIMARY KEY,
    job_id INT NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    job_seeker_id INT NOT NULL REFERENCES job_seeker_profiles(id) ON DELETE CASCADE,
    cover_letter TEXT,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'reviewed', 'interview', 'rejected', 'accepted')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(job_id, job_seeker_id)
);

