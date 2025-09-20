CREATE TABLE IF NOT EXISTS webhook_events (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    event_id VARCHAR(191) NOT NULL,
    event_type VARCHAR(150) NULL,
    paypal_transmission_id VARCHAR(191) NULL,
    payment_id INT NULL,
    paypal_order_id VARCHAR(191) NULL,
    capture_id VARCHAR(191) NULL,
    resource JSON NOT NULL,
    status VARCHAR(50) DEFAULT 'PENDING',
    received_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    processed_at DATETIME NULL,
    UNIQUE KEY uq_event_id (event_id),
    INDEX idx_transmission_id (paypal_transmission_id),
    INDEX idx_event_payment_id (payment_id),
    INDEX idx_event_paypal_order (paypal_order_id),
    CONSTRAINT fk_webhook_payment FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
