import { GraphQLClient } from "graphql-request";

// Knows API location and makes requests
export const graphqlClient = new GraphQLClient(
    // TODO: move to env variable
    "http://localhost:8080/query"
);