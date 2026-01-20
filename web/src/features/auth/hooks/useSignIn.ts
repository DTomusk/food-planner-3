import { graphqlClient, SignInDocument, type SignInMutation, type SignInMutationVariables } from "@/lib";
import { useMutation } from "@tanstack/react-query";
import type { ClientError } from "graphql-request";

export function useSignIn() {
    return useMutation<SignInMutation, ClientError, SignInMutationVariables, unknown>({
        mutationFn: (variables) => graphqlClient.request(SignInDocument, variables),
        onSuccess: (data) => {
            console.log("User signed in, jwt: ", data.signin.jwt);
        },
        onError: (error) => {
            console.error("Error signing in:", error);
        }
    });
}