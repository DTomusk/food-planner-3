import { GetIngredientsDocument, type GetIngredientsQuery, graphqlRequest } from "@/lib";
import { useQuery } from "@tanstack/react-query";
import type { ClientError } from "graphql-request";
import type { IngredientOptionModel } from "../types";

export function useIngredients() {
    return useQuery<GetIngredientsQuery, ClientError, IngredientOptionModel[]>({
        queryKey: ["ingredients"],
        queryFn: async () => graphqlRequest(GetIngredientsDocument),
        select: (data) => data.ingredients.map(ingredient => ({
            id: ingredient.id,
            name: ingredient.name,
            counter: ingredient.counter,
            preferredUnit: {
                val: ingredient.preferredUnit.val,
                name: ingredient.preferredUnit.name,
                symbol: ingredient.preferredUnit.symbol,
            }
        }))
    });
}