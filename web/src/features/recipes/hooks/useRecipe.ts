import { queryOptions, useQuery } from "@tanstack/react-query";
import { GetRecipeDocument, type GetRecipeQuery } from "../../../lib/graphql.generated";
import type { ClientError } from "graphql-request";
import { graphqlRequest } from "@/lib";
import { mapRecipeDetail } from "../mappers/recipeMapper";
import { mapRecipeToFormValues } from "../mappers/recipeFormMapper";
import type { Recipe, User } from "../types";
import type { RecipeFormValues } from "../types";

type UseRecipeResult = {
    recipe: Recipe;
    user: User;
    formValues: RecipeFormValues;
};

export function recipeQueryOptions(id: string) {
    return queryOptions<GetRecipeQuery, ClientError>({
        queryKey: ["recipe", id] as const,
        queryFn: () => graphqlRequest(GetRecipeDocument, { id }),
    });
}

export function useRecipe(id: string) {
    return useQuery<GetRecipeQuery, ClientError, UseRecipeResult>({
        ...recipeQueryOptions(id),
        enabled: Boolean(id),
        select: (data) => {
            if (!data.recipe) {
                throw new Error("Recipe not found");
            }

            return {
                ...mapRecipeDetail(data.recipe),
                formValues: mapRecipeToFormValues(data.recipe),
            };
        },
    });
}