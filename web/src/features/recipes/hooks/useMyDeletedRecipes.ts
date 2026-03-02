import { useQuery } from "@tanstack/react-query";
import { GetMyDeletedRecipesDocument, type GetMyDeletedRecipesQuery } from "../../../lib/graphql.generated";
import { graphqlClient } from "@/lib/graphqlClient";

export function useMyDeletedRecipes() {
    return useQuery<GetMyDeletedRecipesQuery>({
        queryKey: ["recipes","me","deleted"],
        queryFn: () => graphqlClient.request(GetMyDeletedRecipesDocument),
    });
}