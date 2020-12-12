CREATE TABLE tags (
  tag_id SERIAL NOT NULL,
  article_id SERIAL NOT NULL,
  tag text NOT NULL,
  UNIQUE (tag_id),
  PRIMARY KEY (tag_id)
);
