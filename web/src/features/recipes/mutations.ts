import { graphqlClient } from "../../lib/graphqlClient";

const CREATE_RECIPE = `
mutation CreateRecipe($input: NewRecipe!) {
    createRecipe(input: $input) {
        id
        name
    }
}`;

export async function createRecipe(input: { name: string }) {
    return graphqlClient.request(CREATE_RECIPE, { input });
}