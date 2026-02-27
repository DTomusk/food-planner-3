import type { GetRecipeQuery } from "@/lib";
import type { Recipe } from "../types";

export function mapRecipe(
    gqlRecipe: NonNullable<GetRecipeQuery["recipe"]>
): Recipe {
    return {
        id: gqlRecipe.id,
        name: gqlRecipe.name,
        // TODO: add proper formatting later
        ingredients: gqlRecipe.ingredientUsages.map((iu) => ({
            name: iu.ingredient.name,
            formattedQuantity: iu.quantity + iu.unit.symbol
        })),
        user: gqlRecipe.user.id,
        prepMins: gqlRecipe.prepMins,
        cookMins: gqlRecipe.cookMins,
        portions: gqlRecipe.portions,
        source: {
            type: gqlRecipe.source.type === 0 ? "none" : gqlRecipe.source.type === 1 ? "website" : gqlRecipe.source.type === 2 ? "cookbook" : "original",
            url: gqlRecipe.source.url || undefined,
            bookTitle: gqlRecipe.source.bookTitle || undefined,
            bookPage: gqlRecipe.source.bookPage || undefined,
            instructions: gqlRecipe.source.instructions || undefined,
        }
    }
}