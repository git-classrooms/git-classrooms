import { queryOptions } from "@tanstack/react-query";
import { createInfoApi } from "@/lib/utils.ts";

const apiClient = createInfoApi();

export const gitlabInfoQueryOptions = () =>
  queryOptions({
    queryKey: ["gitLabInfo"],
    queryFn: async () => {
      const res = await apiClient.getGitlabInfo();
      return res.data;
    }
  })
