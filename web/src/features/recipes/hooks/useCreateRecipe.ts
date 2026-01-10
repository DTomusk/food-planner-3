import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { CreateRecipeMutation, CreateRecipeMutationVariables } from "../../../lib/graphql.generated";
import type { ClientError } from "graphql-request";
import { recipesApi } from "../api";

export function useCreateRecipe() {
    const queryClient = useQueryClient();

    return useMutation<CreateRecipeMutation, ClientError, CreateRecipeMutationVariables, unknown>({
        mutationFn: (variables) => recipesApi.create( variables.input.name ),
        onSuccess: (data) => {
            queryClient.invalidateQueries({ queryKey: ["recipes"] });
            console.log("Recipe created:", data.createRecipe);
        },
        onError: (error) => {
            console.error("Error creating recipe:", error);
        }
    });
}