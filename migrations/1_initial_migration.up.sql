CREATE TABLE sellers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    url VARCHAR(255) NOT NULL,
    account VARCHAR(255) NOT NULL
);

CREATE TABLE games (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    price DOUBLE NOT NULL,
    on_sale BOOL DEFAULT TRUE,
    seller_id INT NOT NULL,
    FOREIGN KEY (seller_id) REFERENCES sellers (id) ON DELETE CASCADE
);


CREATE TABLE `keys` (
    id INT AUTO_INCREMENT PRIMARY KEY,
    key_string VARCHAR(255) NOT NULL,
    game_id INT NOT NULL,
    seller_id INT NOT NULL,
    status VARCHAR(255) NOT NULL DEFAULT 'available',
    FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE,
    FOREIGN KEY (seller_id) REFERENCES sellers (id) ON DELETE CASCADE
);

CREATE TABLE payment_sessions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    key_id INT NOT NULL,
    price DOUBLE NOT NULL,
    date TIMESTAMP NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    customer_email VARCHAR(255) NOT NUll,
    customer_address VARCHAR(255) NOT NULL,
    FOREIGN KEY (key_id) REFERENCES `keys` (id) ON DELETE CASCADE
);