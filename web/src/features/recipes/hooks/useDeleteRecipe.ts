import { useMutation, useQueryClient } from "@tanstack/react-query";
import { DeleteRecipeDocument, type DeleteRecipeMutation, type DeleteRecipeMutationVariables } from "../../../lib/graphql.generated";
import type { ClientError } from "graphql-request";
import { graphqlClient } from "@/lib/graphqlClient";

export function useDeleteRecipe() {
    const queryClient = useQueryClient();

    return useMutation<DeleteRecipeMutation, ClientError, DeleteRecipeMutationVariables, unknown>({
        mutationFn: (variables) => graphqlClient.request(DeleteRecipeDocument, variables),
        onSuccess: (data) => {
            queryClient.invalidateQueries({ queryKey: ["recipes","me"] });
            console.log("Recipe deleted:", data.deleteRecipe);
        },
        onError: (error) => {
            console.error("Error deleting recipe:", error);
        }
    });
}