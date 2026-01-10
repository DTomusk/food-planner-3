import { useQuery } from "@tanstack/react-query";
import { type GetRecipesQuery } from "../../../lib/graphql.generated";
import { recipesApi } from "../api";

export function useRecipes() {
    return useQuery<GetRecipesQuery>({
        queryKey: ["recipes"],
        queryFn: recipesApi.getAll,
    })
}