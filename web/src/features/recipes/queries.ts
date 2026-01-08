import { graphqlClient } from "../../lib/graphqlClient";

const getRecipes = `
    query GetRecipes {
    recipes {
      name
      id
    }
  }`;

    // The query key is used to identify the query in the cache
  // The query function produces a promise that resolves to the data
export const recipesQuery = {
    // Importantly, query keys are arrays to identify queries in a structured way
    // Two queries with the same key will share the same cache entry
    // You can partially invalidate or refetch queries by their keys
    queryKey: ['recipes'],
    queryFn: () => graphqlClient.request(getRecipes),
}