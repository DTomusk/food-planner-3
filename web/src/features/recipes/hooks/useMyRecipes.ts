import { useQuery } from "@tanstack/react-query";
import { GetMyRecipesDocument, type GetMyRecipesQuery } from "../../../lib/graphql.generated";
import { graphqlClient } from "@/lib/graphqlClient";

export function useMyRecipes() {
    return useQuery<GetMyRecipesQuery>({
        queryKey: ["myRecipes"],
        queryFn: () => graphqlClient.request(GetMyRecipesDocument),
    });
}