USE mycircle;
-- ユーザ1(A-Circle, B-Circle)
-- ユーザ2()
-- ユーザ3(A-Circle)
INSERT INTO users (id, name, email, password_hash, created_at, updated_at) VALUES (1, "ユーザ1", "foo@example.com", "$2a$10$5zIf9lXlK6F7eaMB38uRSes9ecydTeW/xDA53zADvQjrmxA/Q/BsG", NOW(), NOW());
INSERT INTO users (id, name, email, password_hash, created_at, updated_at) VALUES (2, "ユーザ2", "bar@example.com", "$2a$10$5zIf9lXlK6F7eaMB38uRSes9ecydTeW/xDA53zADvQjrmxA/Q/BsG", NOW(), NOW());
INSERT INTO users (id, name, email, password_hash, created_at, updated_at) VALUES (3, "ユーザ3", "baz@example.com", "$2a$10$5zIf9lXlK6F7eaMB38uRSes9ecydTeW/xDA53zADvQjrmxA/Q/BsG", NOW(), NOW());
INSERT INTO circles (id, name, icon_url, created_at, updated_at) VALUES (1, "A-Circle", "http://localhost:4566/my-circle-bucket/circles/54a2e0a21c246a49c4b2f3057ea78da4a38952dbbfa450bc120bde5d99f0a7eb", NOW(), NOW());
INSERT INTO circles (id, name, icon_url, created_at, updated_at) VALUES (2, "B-Circle", "http://localhost:4566/my-circle-bucket/circles/793b589f61a6426a9e6f1891f9ad9db4dfa10b3cea192fe4fa736100e0c02976", NOW(), NOW());
INSERT INTO users_circles (user_id, circle_id, created_at, updated_at) VALUES (1, 1, NOW(), NOW());
INSERT INTO users_circles (user_id, circle_id, created_at, updated_at) VALUES (1, 2, NOW(), NOW());
INSERT INTO users_circles (user_id, circle_id, created_at, updated_at) VALUES (3, 1, NOW(), NOW());