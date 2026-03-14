import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { ClientError } from "graphql-request";
import { graphqlClient } from "@/lib/graphqlClient";
import { UpdateRecipeDocument, type UpdateRecipeMutation, type UpdateRecipeMutationVariables } from "@/lib";

export function useUpdateRecipe() {
    const queryClient = useQueryClient();

    return useMutation<UpdateRecipeMutation, ClientError, UpdateRecipeMutationVariables, unknown>({
        mutationFn: (variables) => graphqlClient.request(UpdateRecipeDocument, variables),
        onSuccess: (data) => {
            queryClient.invalidateQueries({ queryKey: ["recipes"] });
            console.log("Recipe updated:", data.updateRecipe);
        },
        onError: (error) => {
            console.error("Error updating recipe:", error);
        }
    });
}