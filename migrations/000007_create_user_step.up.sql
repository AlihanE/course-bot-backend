create table if not exists user_step(
    id varchar(100) primary key,
    user_id varchar(100) not null,
    step_id varchar(100) not null,
    completed bool,
    sent bool,
    create_date timestamptz,
    update_date timestamptz
);
create index user_steps_idx_user_id on user_step(user_id);
create index user_steps_idx_step_id on user_step(step_id);