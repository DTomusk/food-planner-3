ALTER TABLE users
ADD COLUMN username VARCHAR(255);

UPDATE users 
SET username = CONCAT('user', id);

ALTER TABLE users
ALTER COLUMN username SET NOT NULL UNIQUE;