CREATE DATABASE IF NOT EXISTS ProductsAndUsersDB;
USE ProductsAndUsersDB;

CREATE TABLE IF NOT EXISTS Users (
  id INT PRIMARY KEY,
  name VARCHAR(255),
  mobile VARCHAR(20),
  latitude FLOAT,
  longitude FLOAT,
  created_at DATETIME,
  updated_at DATETIME
);

CREATE TABLE IF NOT EXISTS Products (
  product_id INT PRIMARY KEY AUTO_INCREMENT,
  product_name VARCHAR(255),
  product_description TEXT,
  product_images TEXT,
  product_price DECIMAL(10, 2),
  compressed_product_images TEXT,
  created_at DATETIME,
  updated_at DATETIME
);

INSERT INTO Users (id, name, mobile, latitude, longitude, created_at, updated_at) VALUES
  (1, 'John Doe', '555-1234', 37.7749, -122.4194, '2021-05-01 12:00:00', '2021-05-01 12:00:00'),
  (2, 'Jane Smith', '555-5678', 40.7128, -74.0060, '2021-05-02 09:00:00', '2021-05-03 15:00:00'),
  (3, 'Bob Johnson', '555-9876', 51.5074, -0.1278, '2021-05-03 17:00:00', '2021-05-03 17:00:00');
