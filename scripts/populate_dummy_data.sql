-- Dummy Data Population Script for Career Pulse Backend
-- Run this script after running migrations to populate the database with test data

-- Clear existing data (in reverse order of dependencies)
DELETE FROM messages;
DELETE FROM conversations;
DELETE FROM applications;
DELETE FROM jobs;
DELETE FROM job_seeker_profiles;
DELETE FROM employer_profiles;
DELETE FROM users;

-- Reset sequences
ALTER SEQUENCE users_id_seq RESTART WITH 1;
ALTER SEQUENCE job_seeker_profiles_id_seq RESTART WITH 1;
ALTER SEQUENCE employer_profiles_id_seq RESTART WITH 1;
ALTER SEQUENCE jobs_id_seq RESTART WITH 1;
ALTER SEQUENCE applications_id_seq RESTART WITH 1;
ALTER SEQUENCE conversations_id_seq RESTART WITH 1;
ALTER SEQUENCE messages_id_seq RESTART WITH 1;

-- Insert dummy users
INSERT INTO users (email, password_hash, role, created_at, updated_at) VALUES
-- Job Seekers
('ali.awada@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'job_seeker', NOW(), NOW()),
('sarah.khalil@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'job_seeker', NOW(), NOW()),
('omar.hassan@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'job_seeker', NOW(), NOW()),
('lina.merhi@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'job_seeker', NOW(), NOW()),
('youssef.farhat@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'job_seeker', NOW(), NOW()),

-- Employers
('hr@techcorp.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'employer', NOW(), NOW()),
('jobs@startupx.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'employer', NOW(), NOW()),
('careers@bigtech.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'employer', NOW(), NOW()),
('talent@fintech.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'employer', NOW(), NOW()),
('recruitment@healthtech.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'employer', NOW(), NOW()),

-- Admin
('admin@careerpulse.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', NOW(), NOW());

-- Insert job seeker profiles
INSERT INTO job_seeker_profiles (user_id, first_name, last_name, headline, summary, phone, location, resume_url, logo_url, skills, experience_level, created_at, updated_at) VALUES
(1, 'Ali', 'Awada', '🚀 Full-Stack Developer | Turning Ideas into Scalable Web Apps', 'Passionate full-stack developer with 3+ years of experience building modern web applications. Expert in React, Node.js, and cloud technologies. I love solving complex problems and creating user-friendly solutions.', '+961 81148209', 'Lebanon, Beirut', '/uploads/resumes/ali_resume.pdf', '/uploads/resumes/ali_pfp.jpeg', ARRAY['Python', 'Java', 'Go', 'React', 'Node.js', 'PostgreSQL'], 'Mid-level', NOW(), NOW()),

(2, 'Sarah', 'Khalil', 'Backend Engineer | Go & Microservices Specialist', 'Experienced backend engineer specializing in Go, microservices architecture, and distributed systems. Passionate about clean code, performance optimization, and scalable solutions.', '+961 70123456', 'Lebanon, Tripoli', '/uploads/resumes/sarah_resume.pdf', '/uploads/resumes/sarah_pfp.jpg', ARRAY['Go', 'Docker', 'Kubernetes', 'PostgreSQL', 'Redis', 'gRPC'], 'Senior', NOW(), NOW()),

(3, 'Omar', 'Hassan', 'Frontend Developer | React & TypeScript Expert', 'Creative frontend developer with expertise in React, TypeScript, and modern web technologies. Focused on creating beautiful, responsive, and accessible user interfaces.', '+961 81123456', 'Lebanon, Sidon', '/uploads/resumes/omar_resume.pdf', '/uploads/resumes/omar_pfp.jpg', ARRAY['React', 'TypeScript', 'JavaScript', 'CSS', 'HTML', 'Tailwind'], 'Mid-level', NOW(), NOW()),

(4, 'Lina', 'Merhi', 'Data Scientist | Machine Learning & AI', 'Data scientist with a strong background in machine learning, statistical analysis, and AI. Experienced in Python, R, and various ML frameworks. Passionate about extracting insights from data.', '+961 70123458', 'Lebanon, Byblos', '/uploads/resumes/lina_resume.pdf', '/uploads/resumes/lina_pfp.jpg', ARRAY['Python', 'R', 'TensorFlow', 'PyTorch', 'Pandas', 'NumPy'], 'Senior', NOW(), NOW()),

(5, 'Youssef', 'Farhat', 'DevOps Engineer | Cloud & Infrastructure', 'DevOps engineer specialized in cloud infrastructure, CI/CD pipelines, and automation. Expert in AWS, Docker, Kubernetes, and infrastructure as code.', '+961 81123458', 'Lebanon, Tyre', '/uploads/resumes/youssef_resume.pdf', '/uploads/resumes/youssef_pfp.jpg', ARRAY['AWS', 'Docker', 'Kubernetes', 'Terraform', 'Jenkins', 'Linux'], 'Senior', NOW(), NOW());

-- Insert employer profiles
INSERT INTO employer_profiles (user_id, company_name, industry, website, description, logo_url, location, company_size, created_at, updated_at) VALUES
(6, 'TechCorp Solutions', 'Technology', 'https://techcorp.com', 'Leading technology company specializing in enterprise software solutions. We build innovative products that help businesses scale and succeed in the digital age.', '/uploads/logos/techcorp_logo.png', 'Beirut, Lebanon', '100-500', NOW(), NOW()),

(7, 'StartupX', 'Fintech', 'https://startupx.com', 'Fast-growing fintech startup revolutionizing digital payments and financial services. We are building the future of banking with cutting-edge technology.', '/uploads/logos/startupx_logo.png', 'Beirut, Lebanon', '10-50', NOW(), NOW()),

(8, 'BigTech International', 'Technology', 'https://bigtech.com', 'Global technology giant with offices worldwide. We develop next-generation software solutions and provide cutting-edge technology services to millions of users.', '/uploads/logos/bigtech_logo.png', 'Dubai, UAE', '1000+', NOW(), NOW()),

(9, 'FinTech Innovations', 'Financial Services', 'https://fintech-innovations.com', 'Innovative financial technology company focused on blockchain, cryptocurrency, and digital banking solutions. We are shaping the future of finance.', '/uploads/logos/fintech_logo.png', 'Riyadh, Saudi Arabia', '50-100', NOW(), NOW()),

(10, 'HealthTech Solutions', 'Healthcare', 'https://healthtech.com', 'Healthcare technology company developing AI-powered medical solutions and telemedicine platforms. We are improving healthcare accessibility through technology.', '/uploads/logos/healthtech_logo.png', 'Cairo, Egypt', '100-500', NOW(), NOW());

-- Insert jobs
INSERT INTO jobs (employer_id, title, description, location, job_type, salary_min, salary_max, experience_level, required_skills, category, status, created_at, updated_at) VALUES
(1, 'Senior Backend Developer', 'We are looking for an experienced backend developer to join our engineering team. You will be responsible for designing and implementing scalable APIs, working with microservices architecture, and collaborating with cross-functional teams to deliver high-quality software solutions.', 'Beirut, Lebanon', 'full_time', 2500.00, 4000.00, 'Senior', ARRAY['Go', 'PostgreSQL', 'Docker', 'Microservices'], 'Engineering', 'active', NOW(), NOW()),

(2, 'Frontend React Developer', 'Join our dynamic frontend team to build beautiful and responsive user interfaces. You will work with React, TypeScript, and modern web technologies to create exceptional user experiences for our fintech platform.', 'Remote', 'full_time', 2000.00, 3500.00, 'Mid-level', ARRAY['React', 'TypeScript', 'JavaScript', 'CSS'], 'Engineering', 'active', NOW(), NOW()),

(3, 'DevOps Engineer', 'We need a skilled DevOps engineer to manage our cloud infrastructure and CI/CD pipelines. You will work with AWS, Kubernetes, and automation tools to ensure our systems are scalable, secure, and reliable.', 'Dubai, UAE', 'full_time', 3000.00, 5000.00, 'Senior', ARRAY['AWS', 'Kubernetes', 'Docker', 'Terraform'], 'Engineering', 'active', NOW(), NOW()),

(4, 'Blockchain Developer', 'Exciting opportunity to work on cutting-edge blockchain projects. You will develop smart contracts, DeFi protocols, and blockchain-based applications using Solidity and other blockchain technologies.', 'Riyadh, Saudi Arabia', 'full_time', 2800.00, 4500.00, 'Mid-level', ARRAY['Solidity', 'Web3', 'JavaScript', 'Blockchain'], 'Engineering', 'active', NOW(), NOW()),

(5, 'AI/ML Engineer', 'Join our AI team to develop machine learning models and AI-powered healthcare solutions. You will work with Python, TensorFlow, and medical data to create innovative healthcare technologies.', 'Cairo, Egypt', 'full_time', 2200.00, 3800.00, 'Mid-level', ARRAY['Python', 'TensorFlow', 'Machine Learning', 'AI'], 'Engineering', 'active', NOW(), NOW()),

(1, 'Full Stack Developer', 'We are seeking a versatile full-stack developer to work on our web applications. You will handle both frontend and backend development, working with modern technologies and agile methodologies.', 'Beirut, Lebanon', 'full_time', 1800.00, 3000.00, 'Entry-level', ARRAY['JavaScript', 'Node.js', 'React', 'PostgreSQL'], 'Engineering', 'active', NOW(), NOW()),

(2, 'Product Manager', 'Lead product development for our fintech platform. You will work with engineering teams, stakeholders, and customers to define product requirements and drive product strategy.', 'Remote', 'full_time', 2500.00, 4000.00, 'Senior', ARRAY['Product Management', 'Agile', 'User Research', 'Strategy'], 'Product', 'active', NOW(), NOW()),

(3, 'Data Scientist', 'Analyze large datasets to extract insights and build predictive models. You will work with Python, R, and various ML frameworks to solve complex business problems.', 'Dubai, UAE', 'full_time', 2300.00, 4000.00, 'Mid-level', ARRAY['Python', 'R', 'Machine Learning', 'Statistics'], 'Data Science', 'active', NOW(), NOW());

-- Insert applications
INSERT INTO applications (job_id, job_seeker_id, cover_letter, status, created_at, updated_at) VALUES
(1, 2, 'I am excited to apply for the Senior Backend Developer position at TechCorp Solutions. With my extensive experience in Go and microservices architecture, I believe I would be a valuable addition to your engineering team.', 'pending', NOW(), NOW()),

(2, 3, 'I am very interested in the Frontend React Developer position at StartupX. My experience with React and TypeScript, combined with my passion for creating user-friendly interfaces, makes me a strong candidate for this role.', 'pending', NOW(), NOW()),

(3, 5, 'I would like to apply for the DevOps Engineer position at BigTech International. My expertise in AWS, Kubernetes, and automation tools aligns perfectly with your requirements.', 'pending', NOW(), NOW()),

(4, 1, 'I am applying for the Blockchain Developer position at FinTech Innovations. While my primary experience is in full-stack development, I have been learning blockchain technologies and would love to contribute to your innovative projects.', 'pending', NOW(), NOW()),

(5, 4, 'I am excited to apply for the AI/ML Engineer position at HealthTech Solutions. My background in data science and machine learning, combined with my interest in healthcare technology, makes me an ideal candidate for this role.', 'pending', NOW(), NOW());

-- Insert conversations
INSERT INTO conversations (participant_one_id, participant_two_id, created_at) VALUES
(6, 2, NOW()),
(7, 3, NOW()),
(8, 5, NOW());

-- Insert messages (chat conversations)
INSERT INTO messages (conversation_id, sender_id, content, created_at) VALUES
(1, 6, 'Hi Sarah! Thank you for your application. We were impressed with your experience in Go and microservices. Would you be available for an interview this week?', NOW()),
(1, 2, 'Thank you for considering my application! Yes, I would be happy to schedule an interview. What times work best for you?', NOW()),
(1, 6, 'Great! How about Thursday at 2 PM? We can do it via video call.', NOW()),
(1, 2, 'Perfect! Thursday at 2 PM works for me. I will send you the meeting link.', NOW()),

(2, 7, 'Hi Omar! We reviewed your portfolio and were very impressed with your React projects. Would you like to discuss the Frontend Developer position?', NOW()),
(2, 3, 'Hello! Thank you for reaching out. I would love to discuss the position. When would be a good time to chat?', NOW()),

(3, 8, 'Hi Youssef! Your DevOps experience looks excellent. We would like to invite you for a technical interview to discuss the role in more detail.', NOW()),
(3, 5, 'Hello! Thank you for the invitation. I would be delighted to participate in the technical interview. When would you like to schedule it?', NOW());

-- Display summary
SELECT 'Dummy data populated successfully!' as message;
SELECT 
    'Users: ' || COUNT(*) as users_count 
FROM users;
SELECT 
    'Job Seeker Profiles: ' || COUNT(*) as job_seekers_count 
FROM job_seeker_profiles;
SELECT 
    'Employer Profiles: ' || COUNT(*) as employers_count 
FROM employer_profiles;
SELECT 
    'Jobs: ' || COUNT(*) as jobs_count 
FROM jobs;
SELECT 
    'Applications: ' || COUNT(*) as applications_count 
FROM applications;
SELECT 
    'Conversations: ' || COUNT(*) as conversations_count 
FROM conversations;
SELECT 
    'Messages: ' || COUNT(*) as messages_count 
FROM messages;

