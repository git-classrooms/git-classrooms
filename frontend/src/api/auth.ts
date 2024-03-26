import { apiClient } from "@/lib/utils";
import { queryOptions } from "@tanstack/react-query";

export const authCsrfQueryOptions = queryOptions({
  queryKey: ["csrf_auth"],
  queryFn: async () => {
    const res = await apiClient.get<{ csrf: string }>("/auth/csrf");
    return res.data;
  },
});
