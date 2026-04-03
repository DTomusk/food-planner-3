import type { GetRecipeQuery, GetRecipeVersionQuery, RecipeConnectionListFieldsFragment } from "@/lib";
import type { Recipe, RecipePage, User } from "../types";

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
        version: gqlRecipe.currentVersion.version,
        source: {
            type: gqlRecipe.currentVersion.source.type === 0 ? "none" : gqlRecipe.currentVersion.source.type === 1 ? "website" : gqlRecipe.currentVersion.source.type === 2 ? "cookbook" : "original",
            url: gqlRecipe.currentVersion.source.url || undefined,
            bookTitle: gqlRecipe.currentVersion.source.bookTitle || undefined,
            bookPage: gqlRecipe.currentVersion.source.bookPage || undefined,
            instructions: gqlRecipe.currentVersion.source.instructions || undefined,
        },
        imageUrl: gqlRecipe.currentVersion.imgSrc || null,
        },
        user: {
            id: gqlRecipe.author.id,
            username: gqlRecipe.author.username,
        }
    }
}

export function mapRecipeVersionDetail(
    gqlRecipe: NonNullable<GetRecipeVersionQuery["recipe"]>
): { recipe: Recipe, user: User }  {
    if (!gqlRecipe.version) {
        return {} as any;
    }
    return { recipe:{
        id: gqlRecipe.id,
        name: gqlRecipe.version.name,
        ingredients: gqlRecipe.version.ingredientUsages.map((iu) => ({
            name: iu.ingredient.name,
            quantity: iu.quantity,
            counter: iu.ingredient.counter,
            unitSymbol: iu.unit.symbol,
            plural: iu.ingredient.plural,
            counterPlural: iu.ingredient.counterPlural,
        })),
        prepMins: gqlRecipe.version.prepMins,
        cookMins: gqlRecipe.version.cookMins,
        portions: gqlRecipe.version.portions,
        version: gqlRecipe.version.version,
        source: {
            type: gqlRecipe.version.source.type === 0 ? "none" : gqlRecipe.version.source.type === 1 ? "website" : gqlRecipe.version.source.type === 2 ? "cookbook" : "original",
            url: gqlRecipe.version.source.url || undefined,
            bookTitle: gqlRecipe.version.source.bookTitle || undefined,
            bookPage: gqlRecipe.version.source.bookPage || undefined,
            instructions: gqlRecipe.version.source.instructions || undefined,
        },
        imageUrl: gqlRecipe.version.imgSrc || null,
        },
        user: {
            id: gqlRecipe.author.id,
            username: gqlRecipe.author.username,
        } 
    }
}

export function mapRecipeSummary(
    gqlRecipes: RecipeConnectionListFieldsFragment
): RecipePage {
    return {
        recipes: gqlRecipes.edges.map((edge) => ({
            id: edge.node.id,
            name: edge.node.currentVersion.name,
            imageUrl: "https://images.unsplash.com/photo-1504674900247-0877df9cc836?auto=format&fit=crop&w=800&q=80", // TODO: add image URL to API and use it here
        })),
        endCursor: gqlRecipes.pageInfo.endCursor || null,
        hasNextPage: gqlRecipes.pageInfo.hasNextPage,
    };
}