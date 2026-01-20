import { useQuery } from "@tanstack/react-query";
import { GetRecipeDocument, type GetRecipeQuery } from "../../../lib/graphql.generated";
import type { ClientError } from "graphql-request";
import { graphqlClient } from "@/lib/graphqlClient";

export function useRecipe(id: string) {
    return useQuery<GetRecipeQuery, ClientError, GetRecipeQuery["recipe"]>({
        queryKey: ["recipe", id],
        queryFn: () => graphqlClient.request(GetRecipeDocument, { id }),
        enabled: Boolean(id),
        select: (data) => data.recipe,
    });
}