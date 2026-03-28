import { useInfiniteQuery } from "@tanstack/react-query";
import { GetRecipesDocument, type GetRecipesQuery } from "../../../lib/graphql.generated";
import { graphqlRequest } from "@/lib";
import type { ClientError } from "graphql-request";
import type { RecipePage } from "../types";
import { mergeRecipePages, toRecipePage } from "./recipeConnectionPage";

type UseRecipesParams = {
    first?: number;
    query?: string;
    enabled?: boolean;
};

export function useRecipes(params: UseRecipesParams = {}) {
    const { first = 20, query, enabled = true } = params;
    const normalizedQuery = query?.trim() || null;

    return useInfiniteQuery<GetRecipesQuery, ClientError, RecipePage, readonly ["recipes", number, string | null], string | null>({
        queryKey: ["recipes", first, normalizedQuery],
        enabled,
        initialPageParam: null,
        queryFn: ({ pageParam }) => graphqlRequest(GetRecipesDocument, {
            input: {
                first,
                after: pageParam ?? undefined,
            },
            filter: normalizedQuery ? { query: normalizedQuery } : undefined,
        }),
        getNextPageParam: (lastPage) => {
            if (!lastPage.recipes.pageInfo.hasNextPage) {
                return undefined;
            }
            return lastPage.recipes.pageInfo.endCursor ?? undefined;
        },
        select: (data) => mergeRecipePages(data.pages.map((page) => toRecipePage(page.recipes))),
    });
}