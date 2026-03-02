import type z from "zod";
import type { recipeFormSchema } from "./schemas/recipeFormSchema";

export type Recipe = {
    id: string;
    name: string;
    ingredients: Ingredient[];
    user: string;
    prepMins: number;
    cookMins: number;
    portions: number;
    source: RecipeSource;
};

export type Ingredient = {
    name: string;
    quantity: number;
    counter: string | null | undefined;
    unitSymbol: string;
    plural: string | null | undefined;
    counterPlural: string | null | undefined;
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