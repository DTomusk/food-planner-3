import { useMutation } from "@tanstack/react-query";
import { useAuth } from "./useAuth";
import { graphqlRequest, SignOutDocument, type SignOutMutation } from "@/lib";
import type { ClientError } from "graphql-request";

export function useSignOut() {
    const { signOut } = useAuth();
    return useMutation<SignOutMutation, ClientError>({
        mutationFn: () => graphqlRequest(SignOutDocument),
        onSuccess: () => {
            signOut();
        },
        onError: (error) => {
            console.error("Error signing out:", error);
        }
    });
}