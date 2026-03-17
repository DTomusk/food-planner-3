import { QueryClient } from "@tanstack/query-core";

// Responsible for caching and managing server state
export const queryClient = new QueryClient({
    defaultOptions: {
    queries: {
      retry: 1, 
      retryDelay: attemptIndex => Math.min(1000 * 2 ** attemptIndex, 30000), 
      staleTime: 1000 * 60, 
    },
    mutations: {
      retry: 1,
    },
  },
});