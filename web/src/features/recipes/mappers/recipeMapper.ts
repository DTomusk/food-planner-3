import type { GetRecipeQuery } from "@/lib";
import type { Recipe } from "../types";

export function mapRecipe(
    gqlRecipe: NonNullable<GetRecipeQuery["recipe"]>
): Recipe {
    return {
        id: gqlRecipe.id,
        name: gqlRecipe.name,
        // TODO: add proper formatting later
        ingredients: gqlRecipe.ingredientUsages.map((iu) => iu.quantity + " " + iu.unit.name + " " + iu.ingredient.name),
        user: gqlRecipe.user.id,
    }
}