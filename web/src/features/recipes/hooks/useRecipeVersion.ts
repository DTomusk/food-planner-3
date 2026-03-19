import { useQuery } from "@tanstack/react-query";
import { GetRecipeVersionDocument, type GetRecipeVersionQuery } from "../../../lib/graphql.generated";
import type { ClientError } from "graphql-request";
import { graphqlRequest } from "@/lib";
import { mapRecipeVersionDetail } from "../mappers/recipeMapper";
import type { Recipe, User } from "../types";

type UseRecipeVersionResult = {
    recipe: Recipe;
    user: User;
};

export function useRecipeVersion(id: string, version: number) {
    return useQuery<GetRecipeVersionQuery, ClientError, UseRecipeVersionResult>({
        queryKey: ["recipe", id, "version", version],
        queryFn: () => graphqlRequest(GetRecipeVersionDocument, { id, version }),
        enabled: Boolean(id),
        select: (data) => {
            if (!data.recipe) {
                throw new Error("Recipe not found");
            }

            if (!data.recipe.version) {
                throw new Error("Recipe version not found")
            }

            return {
                ...mapRecipeVersionDetail(data.recipe),
            };
        },
    });
}