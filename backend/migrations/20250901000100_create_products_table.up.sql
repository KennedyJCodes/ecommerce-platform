CREATE TABLE IF NOT EXISTS products (
    product_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL,
    product_description TEXT,
    product_price DECIMAL(10,2) NOT NULL,
    stock_quantity INT NOT NULL DEFAULT 0,
    brand VARCHAR(255) NOT NULL,
    image_url VARCHAR(500),
    movement_type VARCHAR(100)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO products (product_name, product_description, product_price, stock_quantity, brand, image_url, movement_type) VALUES
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
