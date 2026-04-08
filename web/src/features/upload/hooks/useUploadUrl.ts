import { CreateImageUploadUrlDocument, graphqlRequest, type CreateImageUploadUrlMutation, type CreateImageUploadUrlMutationVariables } from "@/lib";
import { useMutation } from "@tanstack/react-query";
import { ClientError } from "graphql-request";

export function useUploadUrl() {
    return useMutation<CreateImageUploadUrlMutation, ClientError, CreateImageUploadUrlMutationVariables, unknown>({
        mutationFn: async (variables) => graphqlRequest(CreateImageUploadUrlDocument, variables),
        onSuccess: (data) => {
            console.log("Upload URL created:", data.createImageUploadUrl);
        },
        onError: (error) => {
            console.error("Error creating upload URL:", error);
        }
    });
}