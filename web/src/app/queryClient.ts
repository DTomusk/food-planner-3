import {
  MutationCache,
  QueryCache,
  QueryClient,
} from "@tanstack/react-query";
import {
  handleUnauthenticatedGraphQLError,
  isUnauthenticatedGraphQLError,
} from "@/lib/auth/unauthenticated";

function shouldRetryRequest(failureCount: number, error: unknown) {
  if (isUnauthenticatedGraphQLError(error)) {
    return false;
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
      retry: shouldRetryRequest,
    },
  },
});