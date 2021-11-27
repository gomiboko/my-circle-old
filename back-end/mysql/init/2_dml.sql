USE mycircle;
INSERT INTO users (name, email, password_hash, created_at, updated_at) VALUES ("ユーザ1", "foo@example.com", "$2a$10$5zIf9lXlK6F7eaMB38uRSes9ecydTeW/xDA53zADvQjrmxA/Q/BsG", NOW(), NOW());
INSERT INTO users (name, email, password_hash, created_at, updated_at) VALUES ("ユーザ2", "bar@example.com", "$2a$10$5zIf9lXlK6F7eaMB38uRSes9ecydTeW/xDA53zADvQjrmxA/Q/BsG", NOW(), NOW());
