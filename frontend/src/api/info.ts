import { queryOptions } from "@tanstack/react-query";
import { createInfoApi } from "@/lib/utils.ts";

export const gitlabInfoQueryOptions = () =>
  queryOptions({
    queryKey: ["gitLabInfo"],
    queryFn: async () => {
      const api = createInfoApi();
      const res = await api.getInfoGitlabResponse();
      return res.data;
    }
  })
