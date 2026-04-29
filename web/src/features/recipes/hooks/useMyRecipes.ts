import { useQuery } from "@tanstack/react-query";
import { GetMyRecipesDocument, type GetMyRecipesQuery } from "../../../lib/graphql.generated";
import { graphqlRequest } from "@/lib";
import type { ClientError } from "graphql-request";
import type { RecipeSummary } from "../types";
import { mapMyRecipeSummaryList } from "../mappers/recipeMapper";

export function useMyRecipes() {
    return useQuery<GetMyRecipesQuery, ClientError, RecipeSummary[]>({
        queryKey: ["me", "recipes"],
        queryFn: () => graphqlRequest(GetMyRecipesDocument),
        select: (data) => {
            if (!data.me?.recipes) {
                throw new Error("No recipes found");
            }

            return mapMyRecipeSummaryList(data.me.recipes);
        },
    });
}