alter table steps
add column step_num int;

create table if not exists reports (
    id varchar(100) primary key,
    user_id varchar(100) not null,
    course_id varchar(100) not null,
    step_id varchar(100) not null,
    text text,
    accepted bool
);
create index reports_idx_user_id on reports(user_id);
create index reports_idx_course_id on reports(course_id);
create index reports_idx_step_id on reports(step_id);

create table if not exists assets (
    id varchar(100) primary key,
    course_id varchar(100) not null,
    step_id varchar(100) not null,
    link text,
    text text,
    picture text
);
create index assets_idx_course_id on assets(course_id);
create index assets_idx_step_id on assets(step_id);

create table if not exists user_courses (
    id varchar(100) primary key,
    user_id varchar(100) not null,
    course_id varchar(100) not null,
    finished bool
);
create index user_courses_idx_course_id on user_courses(course_id);
create index user_courses_idx_user_id on user_courses(user_id);