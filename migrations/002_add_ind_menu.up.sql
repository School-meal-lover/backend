CREATE TABLE "ind_menu_sold" (
    "id" uuid UNIQUE PRIMARY KEY DEFAULT (gen_random_uuid ()),
    "meals_id" uuid REFERENCES "meals" ("id"),
    "sold_out_at" timestamp NULL,
)