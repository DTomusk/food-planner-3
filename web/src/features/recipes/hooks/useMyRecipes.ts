import { useQuery } from "@tanstack/react-query";
import { GetMyRecipesDocument, type GetMyRecipesQuery } from "../../../lib/graphql.generated";
import { graphqlRequest } from "@/lib";
import type { ClientError } from "graphql-request";
import type { RecipeSummary } from "../types";
import { requireRecipeConnection, toRecipeSummaryList } from "./recipeConnectionPage";

export function useMyRecipes() {
    return useQuery<GetMyRecipesQuery, ClientError, RecipeSummary[]>({
        queryKey: ["me", "recipes"],
        queryFn: () => graphqlRequest(GetMyRecipesDocument),
        select: (data) => toRecipeSummaryList(requireRecipeConnection(data.me?.recipes)),
    });
}