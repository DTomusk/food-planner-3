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
        portions: gqlRecipe.portions
    }
}