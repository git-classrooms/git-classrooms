import { Button } from "./ui/button";
import { useCsrf } from "@/provider/csrfProvider";

export const LogoutButton = () => {
  const { csrfToken } = useCsrf();

  return (
    <form method="POST" action="/api/v1/auth/sign-out">
      <input type="hidden" name="csrf_token" value={csrfToken} />
      <Button type="submit" className="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded">
        Logout
      </Button>
    </form>
  );
};
