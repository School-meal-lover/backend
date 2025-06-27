ALTER TABLE menu_items
ADD CONSTRAINT uk_menu_item_per_meal UNIQUE (meals_id, category, name);
