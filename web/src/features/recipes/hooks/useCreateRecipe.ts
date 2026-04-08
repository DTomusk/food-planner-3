import { useMutation, useQueryClient } from "@tanstack/react-query";
import { CreateRecipeDocument, type CreateRecipeMutation, type CreateRecipeMutationVariables } from "../../../lib/graphql.generated";
import type { ClientError } from "graphql-request";
import { graphqlRequest } from "@/lib";

export function useCreateRecipe() {
    const queryClient = useQueryClient();

    return useMutation<CreateRecipeMutation, ClientError, CreateRecipeMutationVariables, unknown>({
        mutationFn: (variables) => graphqlRequest(CreateRecipeDocument, variables),
        onSuccess: (data) => {
            queryClient.invalidateQueries({ queryKey: ["recipes"] });
            queryClient.invalidateQueries({ queryKey: ["me", "recipes"] });
            console.log("Recipe created:", data.createRecipe);
        },
        onError: (error) => {
            console.error("Error creating recipe:", error);
        }
    });
}