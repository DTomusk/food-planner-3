import { graphqlClient, SignUpDocument, type SignUpMutation, type SignUpMutationVariables } from "@/lib";
import { useMutation } from "@tanstack/react-query";
import type { ClientError } from "graphql-request";
import { useAuth } from "./useAuth";

export function useSignUp() {
    const { signIn } = useAuth();
    return useMutation<SignUpMutation, ClientError, SignUpMutationVariables, unknown>({
        mutationFn: (variables) => graphqlClient.request(SignUpDocument, variables),
        onSuccess: (data) => {
            signIn(data.signup.jwt);
        },
        onError: (error) => {
            console.error("Error signing up:", error);
        }
    });
}