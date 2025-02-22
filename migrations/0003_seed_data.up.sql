-- Insert a default user
INSERT INTO users (username, password) 
VALUES 
  ('adminuser', '$2a$10$7n6jlWeU62A7NRyxMVclzuRek62Ar9AYZf6XV4A8b9T.MPYsW8LfG'),
  ('testuser', '$2a$10$0i05/M4YX7ikbxFs6//voO0I5oQ0HqlTR7Zhl6hXUDwe31QyoZmii');

-- Insert some example films referencing the newly inserted users
-- user_id=1 -> 'adminuser'
-- user_id=2 -> 'testuser'

INSERT INTO films (user_id, title, director, release_date, cast, genre, synopsis)
VALUES
  (1, 'First Admin Film', 'Admin Director', '2023-01-01', 'Sample Cast A', 'Action', 'An action-packed admin film.'),
  (2, 'Testuser Film', 'Test Director', '2023-02-01', 'Sample Cast B', 'Drama', 'A dramatic test film.'),
  (2, 'Another Testuser Film', 'Test Director 2', '2023-03-01', 'Sample Cast C', 'Comedy', 'A comedic test film.');