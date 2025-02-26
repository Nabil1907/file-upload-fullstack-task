-- create enum type "gender"
CREATE TYPE "public"."gender" AS ENUM ('MALE', 'FEMALE');
-- create enum type "language"
CREATE TYPE "public"."language" AS ENUM ('en', 'fr', 'es');
-- create enum type "restmethods"
CREATE TYPE "public"."restmethods" AS ENUM ('Get', 'Post');
-- create "students" table
CREATE TABLE "public"."students" (
  "uuid" uuid NOT NULL,
  "student_id" text NULL,
  "student_name" text NULL,
  "subject" text NULL,
  "grade" numeric NULL,
  PRIMARY KEY ("uuid")
);
