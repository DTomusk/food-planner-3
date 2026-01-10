import { GraphQLClient } from "graphql-request";

// Knows API location and makes requests
export const graphqlClient = new GraphQLClient(
    import.meta.env.VITE_API_URL
);