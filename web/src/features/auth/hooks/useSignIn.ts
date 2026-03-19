import { graphqlRequest, SignInDocument, type SignInMutation, type SignInMutationVariables } from "@/lib";
import { useMutation } from "@tanstack/react-query";
import type { ClientError } from "graphql-request";
import { useAuth } from "./useAuth";

export function useSignIn() {
    const { signIn } = useAuth();
    return useMutation<SignInMutation, ClientError, SignInMutationVariables, unknown>({
        mutationFn: (variables) => graphqlRequest(SignInDocument, variables),
        onSuccess: (data) => {
            signIn(data.signin.jwt);
        },
        onError: (error) => {
            console.error("Error signing in:", error);
        }
    });
}