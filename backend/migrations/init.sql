-- ============================================================
-- Initial schema for ecommerce-platform
-- This runs automatically when MySQL container starts for the
-- first time (via docker-entrypoint-initdb.d)
-- ============================================================

-- User accounts table
CREATE TABLE IF NOT EXISTS user_registration (
    UserID    INT            NOT NULL AUTO_INCREMENT,
    UserName  VARCHAR(255)   NOT NULL,
    Password  VARCHAR(255)   NOT NULL,
    PRIMARY KEY (UserID),
    UNIQUE KEY uq_username (UserName)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Products catalog table
CREATE TABLE IF NOT EXISTS Products (
    Product_ID          INT             NOT NULL AUTO_INCREMENT,
    Product_Name        VARCHAR(255)    NOT NULL,
    Product_Description TEXT,
    Product_Price       DECIMAL(10,2)   NOT NULL,
    Stock_Quantity      INT             NOT NULL DEFAULT 0,
    Brand               VARCHAR(255)    NOT NULL,
    Image_URL           VARCHAR(500),
    Movement_Type       VARCHAR(100),
    PRIMARY KEY (Product_ID)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Comments / reviews table
CREATE TABLE IF NOT EXISTS comments (
    ID      INT          NOT NULL AUTO_INCREMENT,
    UserID  INT          NOT NULL,
    Content TEXT         NOT NULL,
    Rating  INT          NOT NULL,
    Date    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (ID),
    CONSTRAINT fk_comments_user FOREIGN KEY (UserID)
        REFERENCES user_registration(UserID)
        ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Payments table (from migration)
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

-- Webhook events table (from migration)
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

-- ============================================================
-- Seed data: sample watches
-- ============================================================
INSERT INTO Products (Product_Name, Product_Description, Product_Price, Stock_Quantity, Brand, Image_URL, Movement_Type) VALUES
('Submariner Date', 'Classic luxury dive watch with ceramic bezel', 12500.00, 15, 'Rolex', '/images/rolex-submariner.jpg', 'Automatic'),
('Speedmaster Moonwatch', 'Legendary chronograph worn on the moon', 7500.00, 20, 'Omega', '/images/omega-speedmaster.jpg', 'Manual'),
('Carrera Chronograph', 'Sporty elegance with precision timing', 5500.00, 12, 'Tag Heuer', '/images/tag-carrera.jpg', 'Automatic'),
('Royal Oak Offshore', 'Iconic octagonal bezel luxury sports watch', 28000.00, 5, 'Audemars Piguet', '/images/ap-royaloak.jpg', 'Automatic'),
('Nautilus 5711', 'Elegant steel sports watch with horizontal embossing', 35000.00, 3, 'Patek Philippe', '/images/patek-nautilus.jpg', 'Automatic'),
('Big Pilot 43', 'Aviator watch with oversized crown and clear dial', 9500.00, 8, 'IWC', '/images/iwc-bigpilot.jpg', 'Automatic'),
('Seamaster Diver 300M', 'Professional dive watch with helium escape valve', 5200.00, 25, 'Omega', '/images/omega-seamaster.jpg', 'Automatic'),
('G-Shock DW5600', 'Tough digital watch with shock resistance', 150.00, 100, 'Casio', '/images/casio-gshock.jpg', 'Quartz'),
('Daytona Cosmograph', 'High-performance racing chronograph', 18000.00, 7, 'Rolex', '/images/rolex-daytona.jpg', 'Automatic'),
('Santos de Cartier', 'Timeless square watch with exposed screws', 8500.00, 10, 'Cartier', '/images/cartier-santos.jpg', 'Automatic'),
('Reverso Classic', 'Art deco reversible watch case', 6500.00, 14, 'Jaeger-LeCoultre', '/images/jlc-reverso.jpg', 'Manual'),
('Tank Louis Cartier', 'Iconic rectangular watch since 1917', 7500.00, 9, 'Cartier', '/images/cartier-tank.jpg', 'Quartz'),
('Navitimer B01', 'Pilot watch with circular slide rule bezel', 8200.00, 11, 'Breitling', '/images/breitling-navitimer.jpg', 'Automatic'),
('Fifty Fathoms', 'First modern dive watch, rich heritage', 12000.00, 6, 'Blancpain', '/images/blancpain-fifty.jpg', 'Automatic'),
('Grand Seiko Spring Drive', 'Japanese precision with smooth sweeping seconds', 6800.00, 18, 'Grand Seiko', '/images/grandseiko-spring.jpg', 'Spring Drive');
