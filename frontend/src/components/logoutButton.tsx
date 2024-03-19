import { Button } from "./ui/button";

export const LogoutButton = () => {
  return (
    <form method="POST" action="/api/v1/auth/sign-out">
      <Button type="submit" className="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded">
        Logout
      </Button>
    </form>
  );
};
