import { ClientError } from "graphql-request";

export function extractErrorMessage(error: unknown): string {
    if (error instanceof ClientError) {
        if (error.response.errors && error.response.errors.length > 0) {
        return error.response.errors.map(e => e.message).join("; ");
        }

        return `GraphQL error: ${error.response.status} ${error.response.statusText || ""}`.trim();
    }

    if (error instanceof Error) {
        return error.message;
    }

    return String(error);
}