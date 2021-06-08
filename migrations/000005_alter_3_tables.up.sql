ALTER TABLE reports
add column create_date timestamptz,
add column update_date timestamptz;

ALTER TABLE assets
add column create_date timestamptz,
add column update_date timestamptz;

ALTER TABLE user_courses
add column create_date timestamptz,
add column update_date timestamptz;
