CREATE TYPE restaurant_type AS ENUM ('RESTAURANT_1', 'RESTAURANT_2');

CREATE TABLE "weeks" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (gen_random_uuid()),
  "start_date" date NOT NULL,
  "restaurant" restaurant_type NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "meals" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (gen_random_uuid()),
  "weeks_id" uuid REFERENCES "weeks"("id"),
  "date" date,
  "day_of_week" varchar NOT NULL,
  "meal_type" varchar NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "menu_items" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (gen_random_uuid()),
  "meals_id" uuid REFERENCES "meals"("id"),
  "category" varchar NOT NULL,
  "name" varchar,
  "name_en" varchar,
  "price" decimal(10, 2),
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

COMMENT ON COLUMN "meals"."day_of_week" IS 'Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday';
COMMENT ON COLUMN "meals"."meal_type" IS 'Breakfast, Lunch_1, Lunch_2, Dinner';
COMMENT ON COLUMN "menu_items"."category" IS '밥, 국, 메인메뉴, 반찬';

ALTER TABLE "menu_items" ADD CONSTRAINT "unique_menu_item_in_meal" UNIQUE ("meals_id", "category", "name");
