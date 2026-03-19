import { GetRecipeVersionsDocument, graphqlRequest, type GetRecipeVersionsQuery } from "@/lib";
import { useQuery } from "@tanstack/react-query";
import type { ClientError } from "graphql-request";

type UseRecipeVersionsResult = {
    version: number;
    createdAt: string;
}

export function useRecipeVersions(id: string) {
    return useQuery<GetRecipeVersionsQuery, ClientError, UseRecipeVersionsResult[]>({
        queryKey: ["recipe", id, "versions"],
        queryFn: () => graphqlRequest(GetRecipeVersionsDocument, { id }),
        enabled: Boolean(id),
        select: (data) => {
            if (!data.recipe) {
                throw new Error("Recipe not found");
            }

            return data.recipe.versions.map((v) => ({
                version: v.version,
                createdAt: v.createdAt,
            }));
        },
    });
}