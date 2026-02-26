export type Recipe = {
    id: string;
    name: string;
    ingredients: Ingredients[];
    user: string;
    prepMins: number;
    cookMins: number;
    portions: number;
};

export type Ingredients = {
    name: string;
    formattedQuantity: string;
}

export type RecipeSummary = {
    id: string;
    name: string;
};

export type RecipeFormValues = {
    name: string;
    ingredientUsages: {
        ingredientId: string;
        quantity: number;
        unit: number;
    }[];
    prepMins: number;
    cookMins: number;
    portions: number;
    sourceType: number;
    url?: string;
    bookTitle?: string;
    bookPage?: number;
    instructions?: string;
};