CREATE TABLE "images" (
                        "id" bigserial PRIMARY KEY,
                         "user_id" bigserial NOT NULL,
                         "image_path" varchar NOT NULL,
                         "image_url" varchar UNIQUE NOT NULL
);

ALTER TABLE "images" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");