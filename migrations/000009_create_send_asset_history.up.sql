create table if not exists asset_history (
    id varchar(100) primary key,
    data text,
    create_date timestamptz
);