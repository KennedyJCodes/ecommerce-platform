CREATE TABLE IF NOT EXISTS comments (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    content TEXT NOT NULL,
    rating INT NOT NULL,
    date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_comments_user_id (user_id),
    CONSTRAINT fk_comments_user FOREIGN KEY (user_id) REFERENCES user_registration(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT chk_comments_rating CHECK (rating BETWEEN 1 AND 5)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
