USE mycircle;
-- ユーザ1(A-Circle, B-Circle)
-- ユーザ2()
-- ユーザ3(A-Circle)
INSERT INTO users (id, name, email, password_hash, created_at, updated_at) VALUES (1, "ユーザ1", "foo@example.com", "$2a$10$5zIf9lXlK6F7eaMB38uRSes9ecydTeW/xDA53zADvQjrmxA/Q/BsG", NOW(), NOW());
INSERT INTO users (id, name, email, password_hash, created_at, updated_at) VALUES (2, "ユーザ2", "bar@example.com", "$2a$10$5zIf9lXlK6F7eaMB38uRSes9ecydTeW/xDA53zADvQjrmxA/Q/BsG", NOW(), NOW());
INSERT INTO users (id, name, email, password_hash, created_at, updated_at) VALUES (3, "ユーザ3", "baz@example.com", "$2a$10$5zIf9lXlK6F7eaMB38uRSes9ecydTeW/xDA53zADvQjrmxA/Q/BsG", NOW(), NOW());
INSERT INTO circles (id, name, created_at, updated_at) VALUES (1, "A-Circle", NOW(), NOW());
INSERT INTO circles (id, name, created_at, updated_at) VALUES (2, "B-Circle", NOW(), NOW());
INSERT INTO users_circles (user_id, circle_id, created_at, updated_at) VALUES (1, 1, NOW(), NOW());
INSERT INTO users_circles (user_id, circle_id, created_at, updated_at) VALUES (1, 2, NOW(), NOW());
INSERT INTO users_circles (user_id, circle_id, created_at, updated_at) VALUES (3, 1, NOW(), NOW());