import { GetUserDocument, graphqlClient } from "@/lib";
import type { GetUserQuery } from "@/lib/graphql.generated";
import { useQuery } from "@tanstack/react-query";
import type { ClientError } from "graphql-request";

export function useUser(id: string) {
    return useQuery<GetUserQuery, ClientError>({
        queryKey: ["user", id],
        queryFn: () => graphqlClient.request(GetUserDocument, { id }),
    });
}