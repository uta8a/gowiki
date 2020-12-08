drop table if exists userinfo;

CREATE TABLE userinfo (
  id SERIAL NOT NULL,
  username varchar(30) NOT NULL,
  password_hash text NOT NULL,
  UNIQUE (id),
  UNIQUE (username),
  PRIMARY KEY (id)
);

-- CREATE TABLE article {
--   id SERIAL NOT NULL,
--   title,
--   tag,
--   content,
--   comment
-- }

/* insert */

INSERT INTO userinfo (username, password_hash) VALUES (
  'test_user',
  'REDACTED1'
);
