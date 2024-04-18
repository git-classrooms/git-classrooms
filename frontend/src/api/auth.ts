import { apiClient } from "@/lib/utils";
import { User } from "@/types/user";
import { queryOptions, useSuspenseQuery } from "@tanstack/react-query";

export const authCsrfQueryOptions = queryOptions({
  queryKey: ["csrf_auth"],
  queryFn: async () => {
    const res = await apiClient.get<{ csrf: string }>("/auth/csrf");
    return res.data;
  },
});

export const useAuth = () =>
  useSuspenseQuery({
    queryKey: ["auth"],
    queryFn: async () => {
      try {
        const res = await apiClient.get<User>("/me");
        return res.data;
      } catch (_) {
        return null;
      }
    },
    retry: false,
    refetchInterval: 10000,
  });
