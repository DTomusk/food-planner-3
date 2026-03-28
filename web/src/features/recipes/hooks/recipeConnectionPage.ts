import type { RecipeConnectionListFieldsFragment } from "@/lib";
import { mapRecipeSummary } from "../mappers/recipeMapper";
import type { RecipePage, RecipeSummary } from "../types";

export function requireRecipeConnection(
    connection: RecipeConnectionListFieldsFragment | null | undefined,
    errorMessage = "No recipes found",
): RecipeConnectionListFieldsFragment {
    if (!connection) {
        throw new Error(errorMessage);
    }

    return connection;
}

export function toRecipePage(connection: RecipeConnectionListFieldsFragment): RecipePage {
    return mapRecipeSummary(connection);
}

export function toRecipeSummaryList(connection: RecipeConnectionListFieldsFragment): RecipeSummary[] {
    return toRecipePage(connection).recipes;
}

export function mergeRecipePages(pages: RecipePage[]): RecipePage {
    const lastPage = pages[pages.length - 1];

    return {
        recipes: pages.flatMap((page) => page.recipes),
        endCursor: lastPage?.endCursor ?? null,
        hasNextPage: lastPage?.hasNextPage ?? false,
    };
}
