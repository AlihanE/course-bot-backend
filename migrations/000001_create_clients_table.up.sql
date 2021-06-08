create table if not exists clients(
    id varchar(100) primary key,
    first_name varchar(200),
    last_name varchar(200),
    login varchar(100),
    chat_id varchar(100),
    active boolean
);