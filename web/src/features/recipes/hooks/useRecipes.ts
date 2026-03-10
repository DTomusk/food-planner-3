import { useQuery } from "@tanstack/react-query";
import { GetRecipesDocument, type GetRecipesQuery } from "../../../lib/graphql.generated";
import { graphqlClient } from "@/lib/graphqlClient";
import { mapRecipeSummary } from "../mappers/recipeMapper";
import type { ClientError } from "graphql-request";
import type { RecipeSummary } from "../types";

export function useRecipes() {
    return useQuery<GetRecipesQuery, ClientError, RecipeSummary[]>({
        queryKey: ["recipes"],
        queryFn: () => graphqlClient.request(GetRecipesDocument),
        select: (data) => {
            if (!data.recipes) {
                throw new Error("No recipes found");
            }
            return mapRecipeSummary(data.recipes);
        },
    });
}