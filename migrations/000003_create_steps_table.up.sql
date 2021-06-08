create table if NOT EXISTS steps(
    id varchar(100) primary key,
    course_id varchar(100),
    name varchar(200),
    type varchar(100) not null,
    file varchar(500),
    text text,
    create_date timestamptz,
    update_date timestamptz
);
create index on steps(course_id);