import type z from "zod";
import type { recipeFormSchema } from "./schemas/recipeFormSchema";

export type Recipe = {
    id: string;
    name: string;
    description: string;
    ingredients: Ingredient[];
    prepMins: number;
    cookMins: number;
    portions: number;
    source: RecipeSource;
    version: number;
    imageUrl: string | null;
    animalProductLevel: number;
    containsGluten: boolean;
    isDraft: boolean;
};

// TODO: move to user feature
export type User = {
    id: string;
    username: string;
}

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

export type RecipePage = {
    recipes: RecipeSummary[];
    endCursor: string | null;
    hasNextPage: boolean;
}

export type RecipeSummary = {
    id: string;
    name: string;
    description: string;
    imageUrl: string | null;
    createdAt: string;
    author: User;
    animalProductLevel: number;
    containsGluten: boolean;
    version: number;
    isDraft: boolean;
};

export type RecipeFormValues = z.infer<typeof recipeFormSchema>;

export const RecipeSourceType = {
  None: "none",
  Website: "website",
  Cookbook: "cookbook",
  Original: "original",
} as const;

export type RecipeSourceType = typeof RecipeSourceType[keyof typeof RecipeSourceType];