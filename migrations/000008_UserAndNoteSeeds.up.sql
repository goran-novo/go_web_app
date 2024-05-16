INSERT INTO users (name, email, password_hash, activated, version) VALUES
('User One', 'user1@example.com', 'hash11111', true, 1),
('User Two', 'user2@example.com', 'hash22222', true, 1),
('User Three', 'user3@example.com', 'hash33333', true, 1),
('User Four', 'user4@example.com', 'hash444444', true, 1),
('User Five', 'user5@example.com', 'hash55555', true, 1),
('User Six', 'user6@example.com', 'hash66666', true, 1),
('User Seven', 'user7@example.com', 'hash77777', true, 1);

DO $$
DECLARE
    user_id bigint;
    lat float;
    lon float;
    note_text text;
BEGIN
    FOR user_id IN SELECT id FROM users LOOP
        FOR i IN 1..(10 + (random() * 10)::int) LOOP
            lat := 45.8150 + (random() - 0.5) / 100.0;  -- Random latitude around Zagreb
            lon := 15.9819 + (random() - 0.5) / 100.0;  -- Random longitude around Zagreb
            note_text := 'Sample note text ' || clock_timestamp();
            
            INSERT INTO notes (user_id, latitude, longitude, text, location)
            VALUES (user_id, lat, lon, note_text, ST_SetSRID(ST_MakePoint(lon, lat), 4326));
        END LOOP;
    END LOOP;
END $$;