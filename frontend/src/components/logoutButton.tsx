import { useSuspenseQuery } from "@tanstack/react-query";
import { Button } from "./ui/button";
import { authCsrfQueryOptions } from "@/api/auth";

export const LogoutButton = () => {
  const { data } = useSuspenseQuery(authCsrfQueryOptions);

  return (
    <form method="POST" action="/api/v1/auth/sign-out">
      <input type="hidden" name="csrf_token" value={data.csrf} />
      <Button type="submit" className="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded">
        Logout
      </Button>
    </form>
  );
};
