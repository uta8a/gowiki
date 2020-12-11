CREATE TABLE group_admins (
  id SERIAL NOT NULL,
  group_name varchar(30) NOT NULL,
  group_admin varchar(30) NOT NULL,
  UNIQUE (id),
  UNIQUE (group_name),
  PRIMARY KEY (id)
);
