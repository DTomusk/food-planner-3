import { useInfiniteQuery } from "@tanstack/react-query";
import { GetRecipesDocument, type GetRecipesQuery } from "../../../lib/graphql.generated";
import { graphqlRequest } from "@/lib";
import { mapRecipeSummary } from "../mappers/recipeMapper";
import type { ClientError } from "graphql-request";
import type { RecipePage } from "../types";

type UseRecipesParams = {
    first?: number;
    query?: string;
};

export function useRecipes(params: UseRecipesParams = {}) {
    const { first = 20, query } = params;
    const normalizedQuery = query?.trim() || null;

    return useInfiniteQuery<GetRecipesQuery, ClientError, RecipePage, readonly ["recipes", number, string | null], string | null>({
        queryKey: ["recipes", first, normalizedQuery],
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
        select: (data) => {
            const mappedPages = data.pages.map((page) => mapRecipeSummary(page.recipes));
            const lastPage = mappedPages[mappedPages.length - 1];

            return {
                recipes: mappedPages.flatMap((page) => page.recipes),
                endCursor: lastPage?.endCursor ?? null,
                hasNextPage: lastPage?.hasNextPage ?? false,
            };
        },
    });
}