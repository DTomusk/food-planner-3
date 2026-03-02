import { useMutation, useQueryClient } from "@tanstack/react-query";
import { UndeleteRecipeDocument, type UndeleteRecipeMutation, type UndeleteRecipeMutationVariables } from "../../../lib/graphql.generated";
import type { ClientError } from "graphql-request";
import { graphqlClient } from "@/lib/graphqlClient";

export function useUndeleteRecipe() {
    const queryClient = useQueryClient();

    return useMutation<UndeleteRecipeMutation, ClientError, UndeleteRecipeMutationVariables, unknown>({
        mutationFn: (variables) => graphqlClient.request(UndeleteRecipeDocument, variables),
        onSuccess: (data) => {
            queryClient.invalidateQueries({ queryKey: ["recipes","me","deleted"] });
            console.log("Recipe undeleted:", data.undeleteRecipe);
        },
        onError: (error) => {
            console.error("Error undeleting recipe:", error);
        }
    });
}