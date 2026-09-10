CREATE TABLE IF NOT EXISTS Components (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- CPU, GPU, Motherboard, RAM, PSU, Cooler, Case
    brand VARCHAR(100),
    model VARCHAR(100),
    socket VARCHAR(50), -- for CPU, Motherboard, Cooler
    form_factor VARCHAR(50), -- for Motherboard, Case, PSU
    wattage INT, -- for PSU (capacity), CPU (TDP), GPU (TDP)
    vram INT, -- for GPU
    ram_type VARCHAR(20), -- DDR4, DDR5
    ram_capacity INT, -- for RAM
    ram_speed INT, -- for RAM
    tier VARCHAR(20), -- For algorithmic weighting (S, A, B, C, D)
    workload_score_gaming INT DEFAULT 0,
    workload_score_editing INT DEFAULT 0,
    workload_score_programming INT DEFAULT 0,
    workload_score_general INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Prices (
    id INT AUTO_INCREMENT PRIMARY KEY,
    component_id INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    store_domain VARCHAR(255),
    is_in_country BOOLEAN DEFAULT TRUE,
    url TEXT,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (component_id) REFERENCES Components(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS CompatibilityRules (
    id INT AUTO_INCREMENT PRIMARY KEY,
    rule_type VARCHAR(100), -- Socket, FormFactor, RAMType, Wattage
    description TEXT
);

-- Basic inserts to test
INSERT INTO Components (name, type, brand, model, socket, tier, workload_score_gaming, workload_score_editing, workload_score_programming, workload_score_general, wattage) VALUES
('AMD Ryzen 5 7600', 'CPU', 'AMD', 'Ryzen 5 7600', 'AM5', 'A', 80, 70, 80, 90, 65),
('Intel Core i5-13400F', 'CPU', 'Intel', 'Core i5-13400F', 'LGA1700', 'B', 75, 75, 75, 85, 65),
('ASUS TUF GAMING B650-PLUS', 'Motherboard', 'ASUS', 'B650-PLUS', 'AM5', 'A', 80, 80, 80, 80, 0),
('NVIDIA RTX 4060', 'GPU', 'NVIDIA', 'RTX 4060', NULL, 'B', 85, 70, 60, 90, 115),
('Corsair RM750e', 'PSU', 'Corsair', 'RM750e', NULL, 'S', 90, 90, 90, 90, 750),
('Corsair Vengeance 32GB DDR5', 'RAM', 'Corsair', 'Vengeance', NULL, 'A', 85, 90, 85, 80, 5),
('Samsung 980 PRO 1TB', 'Storage', 'Samsung', '980 PRO', NULL, 'S', 90, 95, 90, 90, 5),
('NZXT H5 Flow', 'Case', 'NZXT', 'H5 Flow', NULL, 'A', 80, 80, 80, 80, 0),
('Thermalright Peerless Assassin 120', 'Cooler', 'Thermalright', 'PA120', 'AM5', 'S', 90, 90, 90, 90, 0);
