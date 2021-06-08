alter table steps
drop column if exists step_num;

drop table if exists reports;
drop index if exists reports_idx_user_id;
drop index if exists reports_idx_course_id;
drop index if exists reports_idx_step_id;

drop table if exists assets;
drop index if exists assets_idx_course_id;
drop index if exists assets_idx_step_id;

drop table if exists user_courses;
drop index if exists user_courses_idx_course_id;
drop index if exists user_courses_idx_user_id;