CREATE TABLE users (
  id bigint PRIMARY KEY,
  name varchar(256) NOT NULL,
  phone_number varchar(32) NOT NULL,
  country varchar(128) NOT NULL,
  city varchar(128) NOT NULL
)
