CREATE TABLE group_users (
  id SERIAL NOT NULL,
  group_name varchar(30) NOT NULL,
  group_user varchar(30) NOT NULL,
  UNIQUE (id),
  PRIMARY KEY (id)
);
