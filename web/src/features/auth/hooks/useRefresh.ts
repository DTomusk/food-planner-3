import { graphqlClient, RefreshDocument, type RefreshMutation } from "@/lib";
import { useMutation } from "@tanstack/react-query";
import type { ClientError } from "graphql-request";
import { useAuth } from "./useAuth";

export function useRefresh() {
    const { signIn } = useAuth();
    return useMutation<RefreshMutation, ClientError>({
        mutationFn: () => graphqlClient.request(RefreshDocument),
        onError: (error) => {
            console.error("Error refreshing token:", error);
        },
        onSuccess: (data) => {
            signIn(data.refresh.jwt);
        },
    });
}