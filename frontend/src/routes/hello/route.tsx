import {
  createFileRoute,
  Outlet,
  ErrorComponent,
  ErrorComponentProps,
} from "@tanstack/react-router";
import { helloQueryOptions } from "@/api/hello.ts";
import { useSuspenseQuery } from "@tanstack/react-query";
import { AxiosError } from "axios";

export const Route = createFileRoute("/hello")({
  loader: ({ context: { queryClient } }) =>
    queryClient.ensureQueryData(helloQueryOptions),
  component: Hello,
  errorComponent: HelloErrorComponent,
});

export function HelloErrorComponent({ error }: ErrorComponentProps) {
  if (error instanceof AxiosError) {
    if (error.response?.status === 404) {
      return <div>Not found</div>;
    }
  }

  return <ErrorComponent error={error} />;
}

function Hello() {
  const { data } = useSuspenseQuery(helloQueryOptions);

  return (
    <div className="p-2">
      <div>Hello from {data}!</div>
      <div>
        <Outlet />
      </div>
    </div>
  );
}
