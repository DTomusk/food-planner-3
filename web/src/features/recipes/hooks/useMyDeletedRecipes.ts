import { useQuery } from "@tanstack/react-query";
import { GetMyDeletedRecipesDocument, type GetMyDeletedRecipesQuery } from "../../../lib/graphql.generated";
import { graphqlClient } from "@/lib/graphqlClient";

export function useMyDeletedRecipes() {
    return useQuery<GetMyDeletedRecipesQuery>({
        queryKey: ["myDeletedRecipes"],
        queryFn: () => graphqlClient.request(GetMyDeletedRecipesDocument),
    });
}