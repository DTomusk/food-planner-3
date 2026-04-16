import {
  MutationCache,
  QueryCache,
  QueryClient,
} from "@tanstack/react-query";
import { ClientError } from "graphql-request";
import {
  handleUnauthenticatedGraphQLError,
  isUnauthenticatedGraphQLError,
} from "@/lib/auth/unauthenticated";

function shouldRetryRequest(failureCount: number, error: unknown) {
  if (isUnauthenticatedGraphQLError(error)) {
    return false;
  }

  if (error instanceof ClientError) {
    const status = error.response.status;
    const hasGraphQLErrors = (error.response.errors?.length ?? 0) > 0;

    // Do not retry application-level GraphQL failures or client-side HTTP errors.
    if (hasGraphQLErrors || (status >= 400 && status < 500)) {
      return false;
    }
  }

  return failureCount < 1;
}

// Responsible for caching and managing server state
export const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: (error) => {
      handleUnauthenticatedGraphQLError(error);
    },
  }),
  mutationCache: new MutationCache({
    onError: (error) => {
      handleUnauthenticatedGraphQLError(error);
    },
  }),
  defaultOptions: {
    queries: {
      retry: shouldRetryRequest,
      retryDelay: (attemptIndex) =>
        Math.min(1000 * 2 ** attemptIndex, 30000),
      staleTime: 1000 * 60,
    },
    mutations: {
      // Mutations can have side effects, so avoid automatic retries.
      retry: false,
    },
  },
});