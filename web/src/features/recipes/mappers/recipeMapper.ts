import type { GetRecipeQuery, GetRecipesQuery } from "@/lib";
import type { Recipe, RecipeSummary, User } from "../types";

export function mapRecipeDetail(
    gqlRecipe: NonNullable<GetRecipeQuery["recipe"]>
): { recipe: Recipe, user: User }  {
    return { recipe:{
        id: gqlRecipe.id,
        name: gqlRecipe.currentVersion.name,
        ingredients: gqlRecipe.currentVersion.ingredientUsages.map((iu) => ({
            name: iu.ingredient.name,
            quantity: iu.quantity,
            counter: iu.ingredient.counter,
            unitSymbol: iu.unit.symbol,
            plural: iu.ingredient.plural,
            counterPlural: iu.ingredient.counterPlural,
        })),
        prepMins: gqlRecipe.currentVersion.prepMins,
        cookMins: gqlRecipe.currentVersion.cookMins,
        portions: gqlRecipe.currentVersion.portions,
        source: {
            type: gqlRecipe.currentVersion.source.type === 0 ? "none" : gqlRecipe.currentVersion.source.type === 1 ? "website" : gqlRecipe.currentVersion.source.type === 2 ? "cookbook" : "original",
            url: gqlRecipe.currentVersion.source.url || undefined,
            bookTitle: gqlRecipe.currentVersion.source.bookTitle || undefined,
            bookPage: gqlRecipe.currentVersion.source.bookPage || undefined,
            instructions: gqlRecipe.currentVersion.source.instructions || undefined,
        }} ,
        user: {
            id: gqlRecipe.author.id,
            username: gqlRecipe.author.username,
        } 
    }
}

export function mapRecipeSummary(
    gqlRecipes: GetRecipesQuery["recipes"]
): RecipeSummary[] {
    return gqlRecipes.map((gqlRecipe) => ({
        id: gqlRecipe.id,
        name: gqlRecipe.currentVersion.name,
    }));
}