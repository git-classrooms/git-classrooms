import { createAuthApi } from "@/lib/utils";
import { queryOptions, useSuspenseQuery } from "@tanstack/react-query";

export const authCsrfQueryOptions = queryOptions({
  queryKey: ["csrf_auth"],
  queryFn: async () => {
    const api = createAuthApi();
    const res = await api.getCsrf();
    return res.data;
  },
});

export const useAuth = () =>
  useSuspenseQuery({
    queryKey: ["auth"],
    queryFn: async () => {
      try {
        const api = createAuthApi();
        const res = await api.getMe();
        return res.data;
      } catch (_) {
        return null;
      }
    },
    retry: false,
    refetchInterval: 10000,
  });
