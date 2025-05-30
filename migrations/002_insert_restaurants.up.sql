INSERT INTO restaurants (name, name_en) VALUES 
('제 1학생식당', 'Student Cafeteria 1'),
('제 2학생식당', 'Student Cafeteria 2'),
('락락', 'Rakkak Bunsik')
ON CONFLICT (name) DO NOTHING;