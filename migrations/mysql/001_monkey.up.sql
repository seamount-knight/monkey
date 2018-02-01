-- CREATE .up.sql files for migrate scripts
CREATE TABLE IF NOT EXISTS monkey (
  uuid VARCHAR(36) PRIMARY KEY,
  name VARCHAR(100) NOT NULL
);
