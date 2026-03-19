import type { TypedDocumentNode } from "@graphql-typed-document-node/core";
import type { RequestDocument, Variables } from "graphql-request";
import { refreshToken } from "./auth/refresh";
import { isUnauthenticatedGraphQLError } from "./auth/unauthenticated";
import { graphqlClient } from "./graphqlClient";

let refreshInFlight: Promise<string> | null = null;

async function refreshTokenSingleFlight(): Promise<string> {
    if (!refreshInFlight) {
        refreshInFlight = refreshToken().finally(() => {
            refreshInFlight = null;
        });
    }

    return refreshInFlight;
}

export async function graphqlRequest<T = unknown, V extends Variables = Variables>(
    document: RequestDocument | TypedDocumentNode<T, V>,
    variables?: V,
): Promise<T> {
    const executeRequest = () => {
        return graphqlClient.request(
            document as RequestDocument,
            variables as Variables | undefined,
        ) as Promise<T>;
    };

    try {
        return await executeRequest();
    } catch (error) {
        if (!isUnauthenticatedGraphQLError(error)) {
            throw error;
        }

        try {
            await refreshTokenSingleFlight();
        } catch {
            // Preserve the original unauthenticated error so existing sign-out handling still runs.
            throw error;
        }

        return executeRequest();
    }
}