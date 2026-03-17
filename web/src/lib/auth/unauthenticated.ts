import { ClientError } from "graphql-request";

const UNAUTHENTICATED_ERROR_CODE = "UNAUTHENTICATED";

let unauthenticatedHandler: (() => void) | null = null;

export function setUnauthenticatedHandler(handler: (() => void) | null) {
    unauthenticatedHandler = handler;
}

export function isUnauthenticatedGraphQLError(error: unknown): boolean {
    if (!(error instanceof ClientError)) {
        return false;
    }

    return (
        error.response.errors?.some(
            (graphQLError) =>
                graphQLError.extensions?.code === UNAUTHENTICATED_ERROR_CODE,
        ) ?? false
    );
}

export function handleUnauthenticatedGraphQLError(error: unknown) {
    if (!isUnauthenticatedGraphQLError(error)) {
        return;
    }

    unauthenticatedHandler?.();
}