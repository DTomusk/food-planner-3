export type Recipe = {
    id: string;
    name: string;
    ingredients: Ingredients[];
    user: string;
};

export type Ingredients = {
    name: string;
    formattedQuantity: string;
}

export type RecipeSummary = {
    id: string;
    name: string;
};

// Leave unit for now as we only have one, but add later
export type RecipeFormValues = {
    name: string;
    ingredientUsages: {
        ingredientId: string;
        quantity: number;
    }[];
};