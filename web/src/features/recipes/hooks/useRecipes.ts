import { useQuery } from "@tanstack/react-query";
import { GetRecipesDocument, type GetRecipesQuery } from "../../../lib/graphql.generated";
import { graphqlClient } from "@/lib/graphqlClient";

export function useRecipes() {
    return useQuery<GetRecipesQuery>({
        queryKey: ["recipes"],
        queryFn: () => graphqlClient.request(GetRecipesDocument),
    });
}