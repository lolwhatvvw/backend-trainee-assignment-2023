CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "firstname" varchar,
  "lastname" varchar,
  "username" varchar(128) NOT NULL UNIQUE,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "segment" (
  "name" varchar PRIMARY KEY,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "user_segments" (
  "user_id" bigint NOT NULL,
  "segment_name" varchar NOT NULL,
  PRIMARY KEY ("user_id", "segment_name")
);

ALTER TABLE user_segments ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ;

ALTER TABLE user_segments ADD FOREIGN KEY ("segment_name") REFERENCES "segment" ("name") ON DELETE CASCADE;

