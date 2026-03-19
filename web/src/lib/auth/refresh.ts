import { RefreshDocument } from "../graphql.generated";
import { graphqlClient } from "../graphqlClient";
import { setToken } from "./token";
import { isUnauthenticatedGraphQLError } from "./unauthenticated";

export class RefreshError extends Error {
    constructor(message: string, public readonly cause?: unknown) {
        super(message);
        this.name = "RefreshError";
    }
}

export async function refreshToken(): Promise<string> {
    try {
        const data = await graphqlClient.request(RefreshDocument);
        const newToken = data.refresh.jwt;

        if (!newToken) {
            throw new RefreshError("No token returned from refresh");
        }

        setToken(newToken);
        return newToken;
    } catch (error) {
        if (isUnauthenticatedGraphQLError(error)) {
            throw new RefreshError("Unauthenticated error during token refresh", error);
        }
        throw new RefreshError("Failed to refresh token", error);
    }
}