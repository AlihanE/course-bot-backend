ALTER TABLE reports
drop column if exists create_date,
drop column if exists update_date;

ALTER TABLE assets
drop column if exists create_date,
drop column if exists update_date;

ALTER TABLE user_courses
drop column if exists create_date,
drop column if exists update_date;
