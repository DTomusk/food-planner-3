import { useQuery } from "@tanstack/react-query";
import { GetMyRecipesDocument, type GetMyRecipesQuery } from "../../../lib/graphql.generated";
import { graphqlRequest } from "@/lib";
import { mapRecipeSummary } from "../mappers/recipeMapper";
import type { ClientError } from "graphql-request";
import type { RecipeSummary } from "../types";

export function useMyRecipes() {
    return useQuery<GetMyRecipesQuery, ClientError, RecipeSummary[]>({
        queryKey: ["recipes","me"],
        queryFn: () => graphqlRequest(GetMyRecipesDocument),
        select: (data) => {
            if (!data.me?.recipes) {
                throw new Error("No recipes found");
            }
            return mapRecipeSummary(data.me.recipes);
        },
    });
}