import { queryOptions, useQuery } from "@tanstack/react-query";
import type { ClientError } from "graphql-request";
import { GetMeDocument, type GetMeQuery, graphqlRequest } from "@/lib";

type MeUser = {
    id: string;
};

export function meQueryOptions() {
    return queryOptions<GetMeQuery, ClientError>({
        queryKey: ["me"] as const,
        queryFn: () => graphqlRequest(GetMeDocument),
    });
}

export function useMe({ enabled = true }: { enabled?: boolean } = {}) {
    return useQuery<GetMeQuery, ClientError, MeUser>({
        ...meQueryOptions(),
        enabled,
        select: (data) => {
            if (!data.me) {
                throw new Error("User not found");
            }

            return {
                id: data.me.id,
            };
        },
    });
}
