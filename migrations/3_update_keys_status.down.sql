ALTER TABLE `keys`
    DROP COLUMN status;
ALTER TABLE `keys`
    ADD COLUMN status VARCHAR(255) DEFAULT 'available';
