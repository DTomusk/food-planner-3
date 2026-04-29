import { queryOptions, useQuery } from "@tanstack/react-query";
import { GetEditRecipeDocument, type GetEditRecipeQuery } from "../../../lib/graphql.generated";
import type { ClientError } from "graphql-request";
import { graphqlRequest } from "@/lib";
import { mapRecipeToFormValues } from "../mappers/recipeFormMapper";
import type { RecipeFormValues } from "../types";

type UseEditRecipeResult = {
    formValues: RecipeFormValues;
};

export function recipeEditQueryOptions(id: string) {
    return queryOptions<GetEditRecipeQuery, ClientError>({
        queryKey: ["recipe", id, "edit"] as const,
        queryFn: () => graphqlRequest(GetEditRecipeDocument, { id }),
    });
}

export function useEditRecipe(id: string) {
    return useQuery<GetEditRecipeQuery, ClientError, UseEditRecipeResult>({
        ...recipeEditQueryOptions(id),
        enabled: Boolean(id),
        select: (data) => {
            if (!data.recipe) {
                throw new Error("Recipe not found");
            }

            return {
                formValues: mapRecipeToFormValues(data.recipe),
            };
        },
    });
}