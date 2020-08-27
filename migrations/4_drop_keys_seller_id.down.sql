ALTER TABLE `keys`
    ADD COLUMN seller_id int;
ALTER TABLE `keys`
    ADD FOREIGN KEY (seller_id) REFERENCES sellers (id);