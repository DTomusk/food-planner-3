import { graphqlClient, SignUpDocument, type SignUpMutation, type SignUpMutationVariables } from "@/lib";
import { useMutation } from "@tanstack/react-query";
import type { ClientError } from "graphql-request";

export function useSignUp() {
    return useMutation<SignUpMutation, ClientError, SignUpMutationVariables, unknown>({
        mutationFn: (variables) => graphqlClient.request(SignUpDocument, variables),
        onSuccess: (data) => {
            console.log("User signed up, jwt: ", data.signup.jwt);
        },
        onError: (error) => {
            console.error("Error signing up:", error);
        }
    });
}