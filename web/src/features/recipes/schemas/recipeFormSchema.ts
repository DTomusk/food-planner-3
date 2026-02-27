import { z } from "zod";

export const recipeFormSchema = z.object({
    name: z.string().min(1, "Recipe name is required"),
    prepMins: z.number().min(1, "Preparation time must be at least 1 minute"),
    cookMins: z.number().min(1, "Cooking time must be at least 1 minute"),
    portions: z.number().min(1, "Number of portions must be at least 1"),
    ingredientUsages: z.array(z.object({
        ingredientId: z.string().min(1, "Ingredient is required"),
        quantity: z.number().min(0.0001, "Quantity must be greater than 0"),
        unit: z.number()
    })).min(1, "At least one ingredient is required"),
    sourceType: z.string(),
    url: z.string().optional(),
    bookTitle: z.string().optional(),
    bookPage: z.number().optional(),
    instructions: z.string().optional()
});


export const recipeFormSections = {
    details: ["name", "prepMins", "cookMins", "portions"] as const,
    ingredients: ["ingredientUsages"] as const,
    source: ["sourceType", "url", "bookTitle", "bookPage", "instructions"] as const
}