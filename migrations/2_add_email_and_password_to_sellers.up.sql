ALTER TABLE sellers
    ADD COLUMN (
        email VARCHAR(255) NOT NULL UNIQUE,
        encrypted_password VARCHAR(255) NOT NULL
        );