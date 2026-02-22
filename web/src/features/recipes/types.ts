export type Recipe = {
    id: string;
    name: string;
    ingredients: string[];
    user: string;
};

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