DELETE FROM films
  WHERE title IN (
    'First Admin Film',
    'Testuser Film',
    'Another Testuser Film'
  );

DELETE FROM users
  WHERE username IN (
    'adminuser',
    'testuser'
  );
