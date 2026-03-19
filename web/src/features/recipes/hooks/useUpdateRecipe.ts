import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { ClientError } from "graphql-request";
import { graphqlRequest, UpdateRecipeDocument, type UpdateRecipeMutation, type UpdateRecipeMutationVariables } from "@/lib";

export function useUpdateRecipe() {
    const queryClient = useQueryClient();

    return useMutation<UpdateRecipeMutation, ClientError, UpdateRecipeMutationVariables, unknown>({
        mutationFn: (variables) => graphqlRequest(UpdateRecipeDocument, variables),
        onSuccess: (data) => {
            queryClient.invalidateQueries({ queryKey: ["recipes"] });
            console.log("Recipe updated:", data.updateRecipe);
        },
        onError: (error) => {
            console.error("Error updating recipe:", error);
        }
    });
}