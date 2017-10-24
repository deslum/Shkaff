-- MySQL Script generated by MySQL Workbench
-- Ср 25 окт 2017 02:26:52
-- Model: New Model    Version: 1.0
-- MySQL Workbench Forward Engineering


-- -----------------------------------------------------
-- Schema shkaff
-- -----------------------------------------------------
DROP SCHEMA IF EXISTS "shkaff" CASCADE;

-- -----------------------------------------------------
-- Schema shkaff
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS  "shkaff"  ;

-- -----------------------------------------------------
-- Table "shkaff"."users"
-- -----------------------------------------------------
DROP TABLE IF EXISTS "shkaff"."users" CASCADE;

CREATE TABLE  "shkaff"."users" (
  "user_id" SMALLINT NOT NULL,
  "login" VARCHAR(16) NULL,
  "password" VARCHAR(32) NULL,
  "api_token" VARCHAR(32) NULL,
  "first_name" VARCHAR(32) NULL,
  "last_name" VARCHAR(32) NULL,
  "is_active" BOOLEAN NULL,
  "is_admin" BOOLEAN NULL,
  PRIMARY KEY ("user_id"),


-- -----------------------------------------------------
-- Table "shkaff"."types"
-- -----------------------------------------------------
DROP TABLE IF EXISTS "shkaff"."types" CASCADE;

CREATE TABLE  "shkaff"."types" (
  "type_id" TINYINT NOT NULL ,
  "type" VARCHAR(16) NULL,
  "cmd_cli" VARCHAR(16) NULL,
  "cmd_dump" VARCHAR(16) NULL,
  "cmd_restore" VARCHAR(16) NULL,
  PRIMARY KEY ("type_id"),


-- -----------------------------------------------------
-- Table "shkaff"."db_settings"
-- -----------------------------------------------------
DROP TABLE IF EXISTS "shkaff"."db_settings" CASCADE;

CREATE TABLE  "shkaff"."db_settings" (
  "db_id" MEDIUMINT NOT NULL ,
  "custom_name" VARCHAR(32) NULL,
  "server_name" VARCHAR(32) NULL,
  "host" VARCHAR(40) NULL,
  "port" SMALLINT NULL,
  "user_id" SMALLINT NULL,
  "is_active" BOOLEAN NULL,
  "type" TINYINT NOT NULL,
  PRIMARY KEY ("db_id", "type"),
  CONSTRAINT "fk_db_settings_types1"
    FOREIGN KEY ("type")
    REFERENCES "shkaff"."types" ("type_id")
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);


-- -----------------------------------------------------
-- Table "shkaff"."databases"
-- -----------------------------------------------------
DROP TABLE IF EXISTS "shkaff"."databases" CASCADE;

CREATE TABLE  "shkaff"."databases" (
  "db_id" SMALLINT NOT NULL ,
  "server_name" VARCHAR(32) NULL,
  "databases" VARCHAR(45) NULL,
  "create_time" TIMESTAMP NULL,
  "is_active" TINYINT NULL,
  "db_settings_id" MEDIUMINT NOT NULL,
  "db_settings_type" TINYINT NOT NULL,
  PRIMARY KEY ("db_id"),
  CONSTRAINT "fk_databases_db_settings1"
    FOREIGN KEY ("db_settings_id" , "db_settings_type")
    REFERENCES "shkaff"."db_settings" ("db_id" , "type")
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);


-- -----------------------------------------------------
-- Table "shkaff"."tasks"
-- -----------------------------------------------------
DROP TABLE IF EXISTS "shkaff"."tasks" CASCADE;

CREATE TABLE  "shkaff"."tasks" (
  "task_id" SMALLINT NOT NULL ,
  "task_name" VARCHAR(32) NULL,
  "verbose" TINYINT NULL,
  "start_time" TIMESTAMP NULL,
  "is_active" BOOLEAN NULL,
  "thread_count" TINYINT NULL,
  "db_settings_id" MEDIUMINT NOT NULL,
  "db_settings_type" TINYINT NOT NULL,
  PRIMARY KEY ("task_id"),
  CONSTRAINT "fk_tasks_db_settings1"
    FOREIGN KEY ("db_settings_id" , "db_settings_type")
    REFERENCES "shkaff"."db_settings" ("db_id" , "type")
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);


-- -----------------------------------------------------
-- Table "shkaff"."users_has_db_settings"
-- -----------------------------------------------------
DROP TABLE IF EXISTS "shkaff"."users_has_db_settings" CASCADE;

CREATE TABLE  "shkaff"."users_has_db_settings" (
  "users_user_id" SMALLINT NOT NULL,
  "db_settings_db_id" MEDIUMINT NOT NULL,
  PRIMARY KEY ("users_user_id", "db_settings_db_id"),
  CONSTRAINT "fk_users_has_db_settings_users"
    FOREIGN KEY ("users_user_id")
    REFERENCES "shkaff"."users" ("user_id")
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT "fk_users_has_db_settings_db_settings1"
    FOREIGN KEY ("db_settings_db_id")
    REFERENCES "shkaff"."db_settings" ("db_id")
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);


DROP SEQUENCE IF EXISTS "shkaff"."types_type_id_sequence";
CREATE SEQUENCE  "shkaff"."types_type_id_sequence";
ALTER TABLE "shkaff"."types" ALTER COLUMN "type_id" SET DEFAULT NEXTVAL('"shkaff"."types_type_id_sequence"');
DROP SEQUENCE IF EXISTS "shkaff"."db_settings_db_id_sequence";
CREATE SEQUENCE  "shkaff"."db_settings_db_id_sequence";
ALTER TABLE "shkaff"."db_settings" ALTER COLUMN "db_id" SET DEFAULT NEXTVAL('"shkaff"."db_settings_db_id_sequence"');
DROP SEQUENCE IF EXISTS "shkaff"."databases_db_id_sequence";
CREATE SEQUENCE  "shkaff"."databases_db_id_sequence";
ALTER TABLE "shkaff"."databases" ALTER COLUMN "db_id" SET DEFAULT NEXTVAL('"shkaff"."databases_db_id_sequence"');
DROP SEQUENCE IF EXISTS "shkaff"."tasks_task_id_sequence";
CREATE SEQUENCE  "shkaff"."tasks_task_id_sequence";
ALTER TABLE "shkaff"."tasks" ALTER COLUMN "task_id" SET DEFAULT NEXTVAL('"shkaff"."tasks_task_id_sequence"');