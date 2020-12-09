CREATE TABLE users (
  id SERIAL NOT NULL,
  username varchar(30) NOT NULL,
  password_hash text NOT NULL,
  UNIQUE (id),
  UNIQUE (username),
  PRIMARY KEY (id)
);
