-- Ensure guest student user exists for anonymous exam taking
INSERT INTO auth.users (
    id, email, password_hash, role, first_name, last_name, is_active, is_email_verified
) VALUES (
    '00000000-0000-0000-0000-000000000001',
    'guest@memere.et',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'student',
    'Guest',
    'Student',
    true,
    true
) ON CONFLICT (id) DO NOTHING;
