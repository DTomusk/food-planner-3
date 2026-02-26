CREATE TABLE ingredient_usages (
    id UUID PRIMARY KEY,
    recipe_id UUID NOT NULL,
    ingredient_id UUID NOT NULL,
    quantity NUMERIC(10, 2) NOT NULL,
    unit INT NOT NULL,
    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
    FOREIGN KEY (ingredient_id) REFERENCES reference.ingredients(id) ON DELETE CASCADE
);