import { graphqlClient, SignInDocument, type SignInMutation, type SignInMutationVariables } from "@/lib";
import { useMutation } from "@tanstack/react-query";
import type { ClientError } from "graphql-request";
import { useAuth } from "./useAuth";

export function useSignIn() {
    const { signIn } = useAuth();
    return useMutation<SignInMutation, ClientError, SignInMutationVariables, unknown>({
        mutationFn: (variables) => graphqlClient.request(SignInDocument, variables),
        onSuccess: (data) => {
            signIn(data.signin.jwt);
        },
        onError: (error) => {
            console.error("Error signing in:", error);
        }
    });
}