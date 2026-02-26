import { GraphQLClient } from "graphql-request";
import { getToken } from "./auth/token";

export const graphqlClient = new GraphQLClient(
    import.meta.env.VITE_API_URL,
    {
        requestMiddleware: (request) => {
            const token = getToken();
            if (token) {
                request.headers = {
                    ...request.headers,
                    Authorization: `Bearer ${token}`,
                    'Content-Type': 'application/json',
                };
            }
            return request;
        }
    }
);