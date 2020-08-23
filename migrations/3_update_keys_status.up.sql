ALTER TABLE `keys`
    DROP COLUMN status;
ALTER TABLE `keys`
    ADD COLUMN status ENUM('available', 'on_hold', 'sold') DEFAULT 'available';
