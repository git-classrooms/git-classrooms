import { queryOptions } from "@tanstack/react-query";
import axios from "axios";

export const helloQueryOptions = queryOptions({
  queryKey: ["hello"],
  queryFn: async () => {
    const hello = await axios.get<string>("/api/hello");
    return hello.data;
  },
});
