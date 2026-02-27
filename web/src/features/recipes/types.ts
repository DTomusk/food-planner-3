import type z from "zod";
import type { recipeFormSchema } from "./schemas/recipeFormSchema";

export type Recipe = {
    id: string;
    name: string;
    ingredients: Ingredients[];
    user: string;
    prepMins: number;
    cookMins: number;
    portions: number;
    source: RecipeSource;
};

export type Ingredients = {
    name: string;
    formattedQuantity: string;
}

export type RecipeSource = {
    type: RecipeSourceType;
    url?: string;
    bookTitle?: string;
    bookPage?: number;
    instructions?: string;
}

export type RecipeSummary = {
    id: string;
    name: string;
};

export type RecipeFormValues = z.infer<typeof recipeFormSchema>;

export const RecipeSourceType = {
  None: "none",
  Website: "website",
  Cookbook: "cookbook",
  Original: "original",
} as const;

export type RecipeSourceType = typeof RecipeSourceType[keyof typeof RecipeSourceType];