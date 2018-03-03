CREATE SCHEMA IF NOT EXISTS shkaff;

CREATE SEQUENCE shkaff.users_seq;

CREATE TABLE IF NOT EXISTS shkaff.users (
  user_id SMALLINT NOT NULL DEFAULT NEXTVAL ('shkaff.users_seq'),
  login VARCHAR(16) NOT NULL,
  password VARCHAR(32) NOT NULL,
  api_token VARCHAR(32) NOT NULL,
  first_name VARCHAR(32) NOT NULL DEFAULT '',
  last_name VARCHAR(32) NOT NULL DEFAULT '',
  is_active BOOLEAN NOT NULL DEFAULT true,
  is_admin BOOLEAN NOT NULL DEFAULT false,
  is_delete boolean NOT NULL DEFAULT false,
  CONSTRAINT users_id_UNIQUE UNIQUE  (user_id),
  PRIMARY KEY (user_id),
  CONSTRAINT login_UNIQUE UNIQUE  (login),
  CONSTRAINT api_token_UNIQUE UNIQUE  (api_token));

CREATE SEQUENCE shkaff.types_seq;

CREATE TABLE IF NOT EXISTS shkaff.types (
  type_id SMALLINT NOT NULL DEFAULT NEXTVAL ('shkaff.types_seq'),
  type VARCHAR(32) NULL,
  cmd_cli VARCHAR(32) NULL,
  cmd_dump VARCHAR(32) NULL,
  cmd_restore VARCHAR(32) NULL,
  PRIMARY KEY (type_id),
  CONSTRAINT type_id_UNIQUE UNIQUE (type_id),
  CONSTRAINT type_UNIQUE UNIQUE (type));

CREATE SEQUENCE shkaff.db_settings_seq;

CREATE TABLE IF NOT EXISTS shkaff.db_settings (
  db_id SMALLINT NOT NULL DEFAULT NEXTVAL ('shkaff.db_settings_seq'),
  user_id SMALLINT NULL,
  type_id SMALLINT NULL,
  server_name VARCHAR(40) NOT NULL DEFAULT '',
  custom_name VARCHAR(40) NOT NULL DEFAULT '',
  host VARCHAR(40) NOT NULL DEFAULT '0.0.0.0',
  port SMALLINT NULL,
  is_active BOOLEAN NOT NULL DEFAULT true,
  db_user VARCHAR(40) NOT NULL DEFAULT '',
  db_password VARCHAR(40) NOT NULL DEFAULT '',
  is_delete boolean NOT NULL DEFAULT false, 
  PRIMARY KEY (db_id, user_id, type_id),
  CONSTRAINT db_id_UNIQUE UNIQUE  (db_id),
  CONSTRAINT db_name_UNIQUE UNIQUE  (custom_name),
  CONSTRAINT fk_db_settings_types1
    FOREIGN KEY (type_id)
    REFERENCES shkaff.types (type_id)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);
 
 CREATE INDEX fk_db_settings_types1_idx ON shkaff.db_settings (type_id);

CREATE SEQUENCE shkaff.tasks_seq;

CREATE TABLE IF NOT EXISTS shkaff.tasks (
  task_id SMALLINT NOT NULL DEFAULT NEXTVAL ('shkaff.tasks_seq'),
  task_name VARCHAR(32) NOT NULL,
  verb SMALLINT NOT NULL DEFAULT 1,
  is_active boolean NOT NULL DEFAULT false,
  thread_count SMALLINT NULL DEFAULT 4,
  ipv6 BOOLEAN NOT NULL DEFAULT false,
  databases TEXT[] NOT NULL DEFAULT '{}',
  gzip BOOLEAN NOT NULL DEFAULT true,
  db_id SMALLINT NOT NULL,
  months INTEGER[12] NOT NULL DEFAULT '{}',
  dumpfolder VARCHAR(128) NOT NULL DEFAULT '/opt/dump',
  is_delete boolean NOT NULL DEFAULT false,
  days INTEGER[31] NOT NULL DEFAULT '{}',
  hours INTEGER[24] NOT NULL DEFAULT '{}',
  minutes SMALLINT NOT NULL,
  PRIMARY KEY (task_id),
  CONSTRAINT task_id_UNIQUE UNIQUE  (task_id),
  CONSTRAINT task_name_UNIQUE UNIQUE (task_name),
  CONSTRAINT fk_tasks_db_settings1
    FOREIGN KEY (db_id)
    REFERENCES shkaff.db_settings (db_id)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);
 
 CREATE INDEX fk_tasks_db_settings1_idx ON shkaff.tasks (db_id);

INSERT INTO shkaff.users (
    login,
    password,
    api_token)
VALUES (
    'admin',
    MD5('admin'),
    '12345'
);

INSERT INTO shkaff.types (
    cmd_cli,
    cmd_dump,
    cmd_restore,
    "type")
VALUES (
    'mongo',
    'mongodump',
    'mongorestore',
    'mongodb'
);

INSERT INTO shkaff.db_settings (
    user_id,
    type_id,
    port,
    server_name)
VALUES (
    1,
    1,
    27017,
    'TestAdmin'
);

INSERT INTO shkaff.tasks (
    db_id,
    is_active,
    task_name,
    minutes)
VALUES (
    1,
    true,
    'FirstTask',
    20
);
