CREATE TABLE articles (
  article_id SERIAL NOT NULL,
  title text NOT NULL,
  article_path text NOT NULL,
  group_name varchar(30) NOT NULL,
  body text NOT NULL,
  UNIQUE (article_id),
  PRIMARY KEY (article_id)
);
