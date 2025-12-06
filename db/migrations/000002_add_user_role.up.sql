ALTER TABLE users ADD COLUMN role VARCHAR(256) DEFAULT "user" NOT NULL;

INSERT IGNORE INTO users (uuid, username, email, password, role)
VALUES ('e3f779cf-6d4e-4590-8236-7e8a226cb8a4', 'admin', 'admin@admin.com', '$2a$10$DJpFyUcr/2DM80YOQg0g..MtYwyWkfb34vW4g6vbp990r0VyCD11O', 'admin')
