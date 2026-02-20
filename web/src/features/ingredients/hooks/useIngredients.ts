import { GetIngredientsDocument, type GetIngredientsQuery, graphqlClient } from "@/lib";
import { useQuery } from "@tanstack/react-query";

export function useIngredients() {
    return useQuery<GetIngredientsQuery>({
        queryKey: ["ingredients"],
        queryFn: async () => graphqlClient.request(GetIngredientsDocument),
    });
}