import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { ClientError } from "graphql-request";
import { graphqlRequest, UpdateRecipeDocument, type UpdateRecipeMutation, type UpdateRecipeMutationVariables } from "@/lib";

export function useUpdateRecipe() {
    const queryClient = useQueryClient();

    return useMutation<UpdateRecipeMutation, ClientError, UpdateRecipeMutationVariables, unknown>({
        mutationFn: (variables) => graphqlRequest(UpdateRecipeDocument, variables),
        onSuccess: async (data, variables) => {
            const updatedRecipeID = data.updateRecipe.id ?? variables.input.id;

            await Promise.all([
                queryClient.invalidateQueries({ queryKey: ["recipes"] }),
                queryClient.invalidateQueries({ queryKey: ["recipe", updatedRecipeID] }),
                queryClient.invalidateQueries({ queryKey: ["me", "recipes"] }),
            ]);

            console.log("Recipe updated:", data.updateRecipe);
        },
        onError: (error) => {
            console.error("Error updating recipe:", error);
        }
    });
}