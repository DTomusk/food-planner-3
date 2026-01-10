import type { GetRecipeQueryVariables, GetRecipeQuery, GetRecipesQuery, CreateRecipeMutation, CreateRecipeMutationVariables } from "@/lib/graphql.generated"
import { graphqlClient } from "@/lib/graphqlClient"

export const recipesApi = {
    getAll: () => {
        return graphqlClient.request<GetRecipesQuery>(`
            query GetRecipes {
                recipes {
                    id
                    name
                }
            }
        `)
    },

    getById: (id: string) => {
        return graphqlClient.request<GetRecipeQuery, GetRecipeQueryVariables>(
            `
            query GetRecipe($id: ID!) {
                recipe(id: $id) {
                    id
                    name
                }
            }
        `,
            { id }
        )
    },

    create: (name: string) => {
        const variables: CreateRecipeMutationVariables = { input: { name } };
        return graphqlClient.request<CreateRecipeMutation, CreateRecipeMutationVariables>(`
            mutation CreateRecipe($input: NewRecipe!) {
                createRecipe(input: $input) {
                    id
                    name
                }
            }
        `, variables
        );
    }
}