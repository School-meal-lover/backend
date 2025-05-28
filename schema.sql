-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: PostgreSQL
-- Generated at: 2025-05-28T09:34:18.278Z

CREATE TABLE "restaurants" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (gen_random_uuid()),
  "name" varchar NOT NULL,
  "name_en" varchar,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "weeks" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (gen_random_uuid()),
  "start_date" date NOT NULL,
  "restaurants_id" uuid,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "meals" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (gen_random_uuid()),
  "name" varchar,
  "name_en" varchar,
  "weeks_id" uuid,
  "date" date,
  "day_of_week" varchar NOT NULL,
  "meal_type" varchar NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "menu_items" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (gen_random_uuid()),
  "meals_id" uuid,
  "category" varchar NOT NULL,
  "name" varchar,
  "name_en" varchar,
  "price" decimal(10,2),
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

COMMENT ON COLUMN "meals"."day_of_week" IS 'Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday';

COMMENT ON COLUMN "meals"."meal_type" IS 'Breakfast, Lunch_1, Lunch_2, Dinner';

COMMENT ON COLUMN "menu_items"."category" IS '밥, 국, 메인메뉴, 반찬';

ALTER TABLE "weeks" ADD FOREIGN KEY ("restaurants_id") REFERENCES "restaurants" ("id");

ALTER TABLE "meals" ADD FOREIGN KEY ("weeks_id") REFERENCES "weeks" ("id");

ALTER TABLE "menu_items" ADD FOREIGN KEY ("meals_id") REFERENCES "meals" ("id");
