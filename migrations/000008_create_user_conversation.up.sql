create table if not exists client_report_conversation (
    id varchar(100) primary key,
    user_id varchar(100) not null,
    report_id varchar(100) not null,
    text text,
    create_date timestamptz,
    update_date timestamptz
);
create index client_report_conv_idx_user_id on client_report_conversation(user_id);
create index client_report_conv_idx_report_id on client_report_conversation(report_id);