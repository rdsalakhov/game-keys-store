ALTER TABLE games
    ADD COLUMN seller_id INT;
alter table games
    add constraint games_ibfk_1
        foreign key (seller_id) references sellers (id)
            on delete cascade;
create index seller_id
    on games (seller_id);