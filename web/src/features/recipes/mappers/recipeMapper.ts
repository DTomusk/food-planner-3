import type { GetMyRecipesQuery, GetRecipeQuery, GetRecipeVersionQuery, RecipeConnectionListFieldsFragment } from "@/lib";
import type { Recipe, RecipePage, User } from "../types";

export function mapRecipeDetail(
    gqlRecipe: NonNullable<GetRecipeQuery["recipe"]>
): { recipe: Recipe, user: User }  {
    var valuesToUse = gqlRecipe.draftVersion ?? gqlRecipe.currentVersion;
    return { recipe:{
        id: gqlRecipe.id,
        name: valuesToUse.name,
        description: valuesToUse.description,
        ingredients: valuesToUse.ingredientUsages.map((iu) => ({
            name: iu.ingredient.name,
            quantity: iu.quantity,
            counter: iu.ingredient.counter,
            unitSymbol: iu.unit.symbol,
            plural: iu.ingredient.plural,
            counterPlural: iu.ingredient.counterPlural,
        })),
        prepMins: valuesToUse.prepMins,
        cookMins: valuesToUse.cookMins,
        portions: valuesToUse.portions,
        version: valuesToUse.version,
        source: {
            type: valuesToUse.source.type === 0 ? "none" : valuesToUse.source.type === 1 ? "website" : valuesToUse.source.type === 2 ? "cookbook" : "original",
            url: valuesToUse.source.url || undefined,
            bookTitle: valuesToUse.source.bookTitle || undefined,
            bookPage: valuesToUse.source.bookPage || undefined,
            instructions: valuesToUse.source.instructions || undefined,
        },
        imageUrl: valuesToUse.imgSrc || null,
        animalProductLevel: valuesToUse.animalProductLevel,
        containsGluten: valuesToUse.containsGluten,
        isDraft: valuesToUse.publishedAt === null,
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
        description: gqlRecipe.version.description,
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
        animalProductLevel: gqlRecipe.version.animalProductLevel,
        containsGluten: gqlRecipe.version.containsGluten,
        isDraft: gqlRecipe.version.publishedAt === null,
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
            imageUrl: edge.node.currentVersion.imgSrc || null,
            description: edge.node.currentVersion.description,
            createdAt: edge.node.createdAt,
            author: {
                id: edge.node.author.id,
                username: edge.node.author.username,
            },
            animalProductLevel: edge.node.currentVersion.animalProductLevel,
            containsGluten: edge.node.currentVersion.containsGluten,
            version: edge.node.currentVersion.version,
            isDraft: edge.node.currentVersion.publishedAt === null,
        })),
        endCursor: gqlRecipes.pageInfo.endCursor || null,
        hasNextPage: gqlRecipes.pageInfo.hasNextPage,
    };
}

export function mapMyRecipeSummary(
    gqlRecipes: NonNullable<NonNullable<GetMyRecipesQuery["me"]>["recipes"]>
): RecipePage {
    return {
        recipes: gqlRecipes.edges.map((edge) => {
            const displayVersion = edge.node.draftVersion ?? edge.node.currentVersion;

            return {
                id: edge.node.id,
                name: displayVersion.name,
                imageUrl: displayVersion.imgSrc || null,
                description: displayVersion.description,
                createdAt: edge.node.createdAt,
                author: {
                    id: edge.node.author.id,
                    username: edge.node.author.username,
                },
                animalProductLevel: displayVersion.animalProductLevel,
                containsGluten: displayVersion.containsGluten,
                version: displayVersion.version,
                isDraft: displayVersion.publishedAt === null,
            };
        }),
        endCursor: gqlRecipes.pageInfo.endCursor || null,
        hasNextPage: gqlRecipes.pageInfo.hasNextPage,
    };
}

export function mapMyRecipeSummaryList(
    gqlRecipes: NonNullable<NonNullable<GetMyRecipesQuery["me"]>["recipes"]>
) {
    return mapMyRecipeSummary(gqlRecipes).recipes;
}