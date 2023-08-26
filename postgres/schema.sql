CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "firstname" varchar,
  "lastname" varchar,
  "username" varchar(128) UNIQUE,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "segment" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "users_segment" (
  "users_id" bigserial,
  "segment_id" bigserial,
  PRIMARY KEY ("users_id", "segment_id")
);

ALTER TABLE "users_segment" ADD FOREIGN KEY ("users_id") REFERENCES "users" ("id");

ALTER TABLE "users_segment" ADD FOREIGN KEY ("segment_id") REFERENCES "segment" ("id");

