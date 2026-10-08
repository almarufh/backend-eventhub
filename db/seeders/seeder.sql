BEGIN;

INSERT INTO users (email, password) VALUES
('admin@eventhub.com', '$2a$12$eImiTXuWVxfM37uY4JANjOL.8/OHG4jT91RjJ/6Nn4DInG1GgM.iu'),
('organizer1@dev.com', '$2a$12$eImiTXuWVxfM37uY4JANjOL.8/OHG4jT91RjJ/6Nn4DInG1GgM.iu'),
('organizer2@design.com', '$2a$12$eImiTXuWVxfM37uY4JANjOL.8/OHG4jT91RjJ/6Nn4DInG1GgM.iu'),
('john.doe@gmail.com', '$2a$12$eImiTXuWVxfM37uY4JANjOL.8/OHG4jT91RjJ/6Nn4DInG1GgM.iu'),
('jane.smith@yahoo.com', '$2a$12$eImiTXuWVxfM37uY4JANjOL.8/OHG4jT91RjJ/6Nn4DInG1GgM.iu');

INSERT INTO communities (title, description, image, status) VALUES
('Tech Developers Indonesia', 'Komunitas developer software dan pegiat IT se-Indonesia.', 'tech_dev_cover.jpg', 'active'),
('UI/UX Design Club', 'Wadah berbagi ilmu desain antarmuka dan pengalaman pengguna.', 'uiux_club.jpg', 'active'),
('Data Science Circle', 'Diskusi seputar Machine Learning, AI, dan Big Data.', 'data_science.jpg', 'active');

INSERT INTO categories (name) VALUES
('Technology'),
('Design'),
('Data & AI'),
('Career'),
('Networking');

INSERT INTO speakers (name, job, office) VALUES
('Budi Santoso', 'Senior Backend Engineer', 'TechCorp Indonesia'),
('Siti Rahma', 'Principal Product Designer', 'Creative Studio'),
('Ahmad Rizky', 'AI Specialist', 'DataLabs Asia');

INSERT INTO profiles (user_id, name, role, dark_preference, address, job, office, image, description) VALUES
(1, 'Admin EventHub', 'admin', TRUE, 'Jakarta Selatan', 'System Administrator', 'EventHub HQ', 'admin_avatar.png', 'Pengelola sistem utama EventHub.'),
(2, 'Eko Kurniawan', 'organizer', FALSE, 'Bandung', 'Community Manager', 'Tech Dev Indo', 'eko_profile.png', 'Membangun ekosistem developer.'),
(3, 'Nia Wijaya', 'organizer', TRUE, 'Surabaya', 'Lead Designer', 'Design Club', 'nia_profile.png', 'UI/UX Enthusiast.'),
(4, 'John Doe', 'attendee', FALSE, 'Jakarta Barat', 'Junior Developer', 'Startup Local', 'john_avatar.png', 'Selalu ingin belajar hal baru.'),
(5, 'Jane Smith', 'attendee', FALSE, 'Yogyakarta', 'Student', 'Universitas Gadjah Mada', 'jane_avatar.png', 'Mahasiswa Ilmu Komputer.');

INSERT INTO auth_users (user_id, is_active, token, device) VALUES
(1, TRUE, 'token_admin_secret_123456', 'Chrome Windows'),
(2, TRUE, 'token_organizer1_secret_654321', 'Safari macOS'),
(4, TRUE, 'token_john_secret_999888', 'Mobile App Android');

INSERT INTO events (community_id, organizer_id, title, location, description, image, capacity, start_time, end_time) VALUES
(1, 2, 'Golang High Performance Workshop', 'Online via Zoom', 'Belajar optimasi performa backend dengan Go.', 'golang_event.png', 2, NOW() + INTERVAL '3 days', NOW() + INTERVAL '3 days 3 hours'),
(2, 3, 'Mastering Figma for Design Systems', 'Co-Working Space Jakarta', 'Workshop praktis membuat Design System scalable.', 'figma_event.png', 2, NOW() + INTERVAL '7 days', NOW() + INTERVAL '7 days 5 hours'),
(3, 1, 'Introduction to LLMs & RAG Architecture', 'Hybrid (Zoom & Hybrid Studio)', 'Eksplorasi pembuatan aplikasi berbasis AI modern.', 'ai_event.png', 2, NOW() + INTERVAL '10 days', NOW() + INTERVAL '10 days 2 hours');

INSERT INTO notifications (title, message, event_id) VALUES
('Pendaftaran Berhasil', 'Terima kasih telah mendaftar di workshop Golang.', 1),
('Perubahan Jadwal', 'Event Mastering Figma akan dimulai 15 menit lebih awal.', 2),
('Pengingat Event', 'Event AI akan dilaksanakan besok jam 19.00 WIB.', 3);

INSERT INTO discussions (user_id, event_id, community_id, message) VALUES
(4, 1, 1, 'Halo, apakah slide materi akan dibagikan setelah event selesai?'),
(2, 1, 1, 'Tentu, slide dan rekaman akan dikirim lewat email.'),
(5, 2, 2, 'Apakah ada prasyarat software yang harus di-install sebelum ikut?');

INSERT INTO user_notifications (user_id, notification_id, is_read, read_at) VALUES
(4, 1, TRUE, NOW() - INTERVAL '1 hour'),
(5, 2, FALSE, NULL),
(4, 3, FALSE, NULL);

INSERT INTO joined_events_users (user_id, event_id) VALUES
(4, 1),
(5, 1),
(4, 2),
(5, 3);

INSERT INTO communities_users (user_id, community_id) VALUES
(2, 1),
(4, 1),
(3, 2),
(5, 2),
(1, 3),
(4, 3);

INSERT INTO events_categories (category_id, event_id) VALUES
(1, 1), 
(5, 1), 
(2, 2), 
(1, 3), 
(3, 3); 

INSERT INTO communities_categories (category_id, community_id) VALUES
(1, 1), 
(2, 2), 
(3, 3); 

INSERT INTO speakers_events (speaker_id, event_id) VALUES
(1, 1),
(2, 2),
(3, 3);

COMMIT;