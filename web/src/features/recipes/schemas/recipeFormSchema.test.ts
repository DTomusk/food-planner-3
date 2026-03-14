import { describe, expect, it } from "vitest";
import { z } from "zod";
import { RecipeSourceType } from "../types";
import { recipeFormSchema } from "./recipeFormSchema";

const baseFormValues = {
    name: "Pasta",
    prepMins: 10,
    cookMins: 15,
    portions: 2,
    ingredientUsages: [{ ingredientId: "ingredient-1", quantity: 1, unit: 1 }],
    sourceType: RecipeSourceType.None,
    url: "",
    bookTitle: "",
    bookPage: undefined,
    instructions: "",
};

describe("recipeFormSchema", () => {
    it("requires a valid URL when website source is selected", () => {
        const result = recipeFormSchema.safeParse({
            ...baseFormValues,
            sourceType: RecipeSourceType.Website,
            url: "not-a-url",
        });

        expect(result.success).toBe(false);

        if (!result.success) {
            expect(z.flattenError(result.error).fieldErrors.url).toContain("Enter a valid URL");
        }
    });

    it("ignores url when another source type is selected", () => {
        const result = recipeFormSchema.safeParse({
            ...baseFormValues,
            sourceType: RecipeSourceType.None,
            url: "not-a-url",
        });

        expect(result.success).toBe(true);
    });

    it("requires cookbook title and page when cookbook source is selected", () => {
        const result = recipeFormSchema.safeParse({
            ...baseFormValues,
            sourceType: RecipeSourceType.Cookbook,
            bookTitle: "",
            bookPage: undefined,
        });

        expect(result.success).toBe(false);

        if (!result.success) {
            const errors = z.flattenError(result.error).fieldErrors;
            expect(errors.bookTitle).toContain("Cookbook title is required");
            expect(errors.bookPage).toContain("Page number is required");
        }
    });

    it("requires instructions only when original source is selected", () => {
        const result = recipeFormSchema.safeParse({
            ...baseFormValues,
            sourceType: RecipeSourceType.Original,
            instructions: "",
        });

        expect(result.success).toBe(false);

        if (!result.success) {
            expect(z.flattenError(result.error).fieldErrors.instructions).toContain("Instructions are required");
        }
    });
});