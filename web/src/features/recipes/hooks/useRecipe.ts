import { useQuery } from "@tanstack/react-query";
import type { GetRecipeQuery } from "../../../lib/graphql.generated";
import { recipesApi } from "../api";
import type { ClientError } from "graphql-request";

export function useRecipe(id: string) {
    return useQuery<GetRecipeQuery, ClientError, GetRecipeQuery["recipe"]>({
        queryKey: ["recipe", id],
        queryFn: () => recipesApi.getById(id),
        enabled: Boolean(id),
        select: (data) => data.recipe,
    });
}