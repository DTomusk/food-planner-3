import { useMutation, useQueryClient } from "@tanstack/react-query";
import { CreateRecipeDocument, type CreateRecipeMutation, type CreateRecipeMutationVariables } from "../../../lib/graphql.generated";
import type { ClientError } from "graphql-request";
import { graphqlClient } from "@/lib/graphqlClient";

export function useCreateRecipe() {
    const queryClient = useQueryClient();

    return useMutation<CreateRecipeMutation, ClientError, CreateRecipeMutationVariables, unknown>({
        mutationFn: (variables) => graphqlClient.request(CreateRecipeDocument, variables),
        onSuccess: (data) => {
            queryClient.invalidateQueries({ queryKey: ["recipes"] });
            console.log("Recipe created:", data.createRecipe);
        },
        onError: (error) => {
            console.error("Error creating recipe:", error);
        }
    });
}