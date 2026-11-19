-- =========================================================
-- DATABASE PRE-LOVED PRODUCT MICROSERVICE (preloved_product)
-- Sesuai spesifikasi resmi Postman Tim Produk
-- =========================================================

CREATE DATABASE IF NOT EXISTS preloved_product
CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;

USE preloved_product;

-- 1. Tabel Kategori Produk
CREATE TABLE IF NOT EXISTS categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    slug VARCHAR(150) NOT NULL UNIQUE,
    description TEXT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- 2. Tabel Produk Pre-Loved
CREATE TABLE IF NOT EXISTS products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    seller_id INT NOT NULL DEFAULT 1,
    category_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NULL,
    price DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    original_price DECIMAL(15, 2) NULL,
    `condition` VARCHAR(50) NOT NULL DEFAULT 'lightly_used', 
    status VARCHAR(30) NOT NULL DEFAULT 'available', -- 'available', 'reserved', 'sold'
    stock INT NOT NULL DEFAULT 1,
    weight_grams INT NOT NULL DEFAULT 500,
    primary_image VARCHAR(500) NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    INDEX idx_seller_id (seller_id),
    INDEX idx_category_id (category_id),
    INDEX idx_status (status)
) ENGINE=InnoDB;

-- 3. Tabel Galeri Foto Produk
CREATE TABLE IF NOT EXISTS product_images (
    id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT NOT NULL,
    image_url VARCHAR(500) NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    display_order INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    INDEX idx_product_id (product_id)
) ENGINE=InnoDB;

-- =========================================================
-- SEED DATA AWAL (Sesuai Postman Tim Produk)
-- =========================================================

-- Seed Kategori
INSERT INTO categories (id, name, slug, description) VALUES
(1, 'Fashion & Pakaian', 'fashion-pakaian', 'Pakaian, jaket, dan busana pre-loved pria, wanita, serta anak'),
(2, 'Sepatu & Sneakers', 'sepatu-sneakers', 'Sepatu kasual, sneakers edisi terbatas, boots, dan alas kaki pre-loved'),
(3, 'Kamera & Fotografi Vintage', 'kamera-vintage', 'Kamera analog, lensa manual, dan aksesori fotografi pre-loved')
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- Seed Produk Pre-Loved
INSERT INTO products (id, seller_id, category_id, name, slug, description, price, original_price, `condition`, status, stock, weight_grams, primary_image) VALUES
(1, 1, 1, "Jaket Denim Levi's Vintage 90s Size L", "jaket-denim-levis-vintage-90s-size-l-abc123", "Kondisi 9/10, warna denim masih pekat, tidak ada sobek, kancing original lengkap.", 450000.00, 1200000.00, 'lightly_used', 'available', 1, 850, "https://images.unsplash.com/photo-1576995853123-5a10305d93c0?w=400&auto=format&fit=crop&q=80"),
(2, 1, 2, "Sepatu Sneakers Converse 70s Chuck Taylor", "sepatu-sneakers-converse-70s-chuck-taylor-def456", "Ukuran 42, sol masih tebal, terawat rapi siap pakai.", 350000.00, 999000.00, 'like_new', 'available', 1, 950, "https://images.unsplash.com/photo-1525966222134-fcfa99b8ae77?w=400&auto=format&fit=crop&q=80"),
(3, 1, 3, "Kamera Analog Vintage 35mm Terawat", "kamera-analog-vintage-35mm-terawat-ghi789", "Kondisi mekanik 100% normal, lensa jernih bebas jamur, viewfinder bersih.", 750000.00, 2500000.00, 'well_used', 'available', 1, 650, "https://images.unsplash.com/photo-1526170375885-4d8ecf77b99f?w=400&auto=format&fit=crop&q=80")
ON DUPLICATE KEY UPDATE name=VALUES(name), price=VALUES(price), status=VALUES(status), stock=VALUES(stock), primary_image=VALUES(primary_image);

-- Seed Galeri Foto
INSERT INTO product_images (product_id, image_url, is_primary, display_order) VALUES
(1, "https://images.unsplash.com/photo-1576995853123-5a10305d93c0?w=400&auto=format&fit=crop&q=80", TRUE, 0),
(2, "https://images.unsplash.com/photo-1525966222134-fcfa99b8ae77?w=400&auto=format&fit=crop&q=80", TRUE, 0),
(3, "https://images.unsplash.com/photo-1526170375885-4d8ecf77b99f?w=400&auto=format&fit=crop&q=80", TRUE, 0);
