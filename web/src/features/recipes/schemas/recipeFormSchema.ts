import { z } from "zod";
import { RecipeSourceType } from "../types";

const sourceTypeOptions = [
    RecipeSourceType.None,
    RecipeSourceType.Website,
    RecipeSourceType.Cookbook,
    RecipeSourceType.Original,
] as const;

const urlSchema = z.url("Enter a valid URL");

export const recipeFormSchema = z.object({
    name: z.string().min(1, "Recipe name is required").max(100, "Recipe name must be at most 100 characters"),
    description: z.string().trim().max(200, "Description must be at most 200 characters").optional(),
    prepMins: z.number().min(1, "Preparation time must be at least 1 minute"),
    cookMins: z.number().min(1, "Cooking time must be at least 1 minute"),
    portions: z.number().min(1, "Number of portions must be at least 1"),
    ingredientUsages: z.array(z.object({
        ingredientId: z.string().min(1, "Ingredient is required"),
        quantity: z.number().min(0.0001, "Quantity must be greater than 0"),
        unit: z.number()
    })).min(1, "At least one ingredient is required"),
    sourceType: z.enum(sourceTypeOptions),
    url: z.string().trim().optional(),
    bookTitle: z.string().trim().optional(),
    bookPage: z.number().int("Page number must be a whole number").min(1, "Page number must be at least 1").optional(),
    instructions: z.string().trim().optional(),
    imgSrc: z.string().trim().optional(),
}).superRefine((values, ctx) => {
    switch (values.sourceType) {
        case RecipeSourceType.Website:
            if (!values.url) {
                ctx.addIssue({
                    code: "custom",
                    message: "URL is required",
                    path: ["url"],
                });
            } else if (!urlSchema.safeParse(values.url).success) {
                ctx.addIssue({
                    code: "custom",
                    message: "Enter a valid URL",
                    path: ["url"],
                });
            }
            break;
        case RecipeSourceType.Cookbook:
            if (!values.bookTitle) {
                ctx.addIssue({
                    code: "custom",
                    message: "Cookbook title is required",
                    path: ["bookTitle"],
                });
            }

            if (values.bookPage === undefined) {
                ctx.addIssue({
                    code: "custom",
                    message: "Page number is required",
                    path: ["bookPage"],
                });
            }
            break;
        case RecipeSourceType.Original:
            if (!values.instructions) {
                ctx.addIssue({
                    code: "custom",
                    message: "Instructions are required",
                    path: ["instructions"],
                });
            }
            break;
        default:
            break;
    }
});


export const recipeFormSections = {
    details: ["name", "prepMins", "cookMins", "portions", "imgSrc"] as const,
    ingredients: ["ingredientUsages"] as const,
    source: ["sourceType", "url", "bookTitle", "bookPage", "instructions"] as const
}