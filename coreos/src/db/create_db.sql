CREATE DATABASE IF NOT EXISTS monkeydb DEFAULT CHARSET=utf8;

USE monkeydb;
/* This table contains our monkeys. */
DROP TABLE IF EXISTS monkeys;
CREATE TABLE monkeys (
  monkeyId BIGINT NOT NULL AUTO_INCREMENT,
  monkeyName varchar(256) DEFAULT NULL,
	/* birthDate is in seconds since UNIX Epoch */
	birthDate INTEGER UNSIGNED DEFAULT NULL,
  PRIMARY KEY (monkeyId)
);

/* Insert some data into monkeys. */
LOCK TABLES monkeys WRITE;
INSERT INTO monkeys (monkeyName, birthDate) VALUES
  ('Janelle', UNIX_TIMESTAMP('2009-03-15 14:01:43')),
  ('Billy', UNIX_TIMESTAMP('2008-01-15 11:00:15'));
UNLOCK TABLES;
