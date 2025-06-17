-- seed_data.sql

-- Users
INSERT INTO User (email, user_name, password, first_name, last_name, display_name, is_active)
VALUES
  ('estecat@yahoo.com', 'estecat', 'password', 'Esteban', 'Price', 'Bubba', 0),
  ('whispers_from_wallphace@gmail.com', 'wallphace', 'password', 'Marcus', 'Price', 'Whispers From Wallphace', 0),
  ('d.cooper@fbi.gov', 'dalecooper', 'password', 'Dale', 'Cooper', 'Coffee Fre@k', 0),
  ('audrey@hornesdepartmentstore.com', 'audrey', 'password', 'Audrey', 'Horne', 'Audrey', 0),
  ('bobby.briggs@twinpeakswa.gov', 'bobbybriggs', 'password', 'Bobby', 'Briggs', 'Bobby', 0),
  ('donna.hayward@twinpeaksclinic.com', 'donnahayward', 'password', 'Donna', 'Hayward', 'Donna', 0);

-- Posts (assumes the above inserts generate ids 1–6 in the same order)
INSERT INTO Post (user_id, content, image)
VALUES
  (3, 'Diane! I''m holding in my hand a small box of chocolate bunnies.', 'chocolate-bunnies.png'),
  (2, 'Esteban knocked over the plant again. 3rd time this week.', ''),
  (2, '', 'modular_patch.png'),
  (3, 'Nothing beats a damn fine cup of coffee in the morning.', 'double-r-diner.jpg'),
  (4, 'High school sucks. At least I’ve got my bike.', ''),
  (5, 'Sometimes I think I can still hear her... Laura.', 'laura-locket.jpg'),
  (2, 'Recording whispers again tonight. Got a new contact mic.', ''),
  (1, '', 'esteban-snoozing.jpg'),
  (5, 'The woods have secrets. I’m starting to believe it now.', ''),
  (4, 'Cruised past the sheriff''s station. Coop’s still parked outside.', 'bobby-bike.jpg');
