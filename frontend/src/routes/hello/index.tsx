import { createFileRoute } from "@tanstack/react-router";
import { helloQueryOptions } from "@/api/hello.ts";

export const Route = createFileRoute("/hello/")({
  loader: ({ context: { queryClient } }) =>
    queryClient.ensureQueryData(helloQueryOptions),
  component: Index,
});

function Index() {
  const data = Route.useLoaderData();

  return (
    <div className="p-2">
      <h3>Welcome Home!!! {data}</h3>
    </div>
  );
}
