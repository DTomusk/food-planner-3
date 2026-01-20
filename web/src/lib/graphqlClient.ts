import { GraphQLClient } from "graphql-request";
import { getToken } from "./auth/token";

export const graphqlClient = new GraphQLClient(
    import.meta.env.VITE_API_URL,
    {
        headers: {
            Authorization: `Bearer ${getToken()}`,
        }
    }
);