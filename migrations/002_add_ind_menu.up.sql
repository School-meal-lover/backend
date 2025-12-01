CREATE TABLE "ind_menu_sold" (
    "id" uuid UNIQUE PRIMARY KEY DEFAULT (gen_random_uuid ()),
    "meals_id" uuid NOT NULL REFERENCES "meals" ("id"),
    "sold_out_at" timestamp NULL,
    CONSTRAINT ind_menu_sold_meals_id_key UNIQUE ("meals_id")
)