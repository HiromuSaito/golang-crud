CREATE TABLE IF NOT EXISTS users (
  id        int     PRIMARY KEY                comment 'ユーザーID'
 ,name     varchar(200)   DEFAULT NULL UNIQUE  comment '氏名'
);