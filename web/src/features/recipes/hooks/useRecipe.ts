import { useQuery } from "@tanstack/react-query";
import { GetRecipeDocument, type GetRecipeQuery } from "../../../lib/graphql.generated";
import type { ClientError } from "graphql-request";
import { graphqlClient } from "@/lib/graphqlClient";
import { mapRecipeDetail } from "../mappers/recipeMapper";
import type { Recipe, User } from "../types";

export function useRecipe(id: string) {
    return useQuery<GetRecipeQuery, ClientError, { recipe: Recipe, user: User }>({
        queryKey: ["recipe", id],
        queryFn: () => graphqlClient.request(GetRecipeDocument, { id }),
        enabled: Boolean(id),
        select: (data) => {
            if (!data.recipe) {
                throw new Error("Recipe not found");
            }

            return mapRecipeDetail(data.recipe);},
    });
}