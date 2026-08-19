CREATE TABLE IF NOT EXISTS payments (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    paypal_order_id VARCHAR(191) NOT NULL,
    capture_id VARCHAR(191),
    status VARCHAR(50) NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    user_id INT NULL,
    raw_response JSON NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_paypal_order_id (paypal_order_id),
    INDEX idx_payments_user_id (user_id),
    CONSTRAINT fk_payments_user FOREIGN KEY (user_id) REFERENCES user_registration(UserID) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
