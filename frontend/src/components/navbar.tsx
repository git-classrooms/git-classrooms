import GitlabLogo from "@/assets/gitlab_logo.svg";
import { ModeToggle } from "@/components/modeToggle.tsx";
import { Link } from "@tanstack/react-router";
import { useCsrf } from "@/provider/csrfProvider.tsx";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { LogOut, Settings, User as UserIcon } from "lucide-react";
import { Avatar } from "./avatar";
import { GetMe } from "@/swagger-client";

export function Navbar(props: { auth: GetMe | null }) {
  const { csrfToken } = useCsrf();

  return (
    <nav className="flex justify-between px-8 py-2.5 mb-8 border-b">
      <div className="flex items-center">
        <a href="/" className="">
          <img className="h-14" src={GitlabLogo} alt="Gitlab Logo" />
        </a>
        <ul className="flex">
          <li className="content-center">
            <Link
              to="/"
              className="font-medium text-sm px-4 py-2 hover:underline"
              activeProps={{ className: "!font-bold" }}
            >
              Dashboard
            </Link>
          </li>
          <li className="content-center">
            <Link
              to="/"
              className="font-medium text-sm px-4 py-2 hover:underline"
              activeProps={{ className: "!font-bold" }}
            >
              Created Classrooms
            </Link>
          </li>
          <li className="content-center">
            <Link
              to="/"
              className="font-medium text-sm px-4 py-2 hover:underline"
              activeProps={{ className: "!font-bold" }}
            >
              Joined Classrooms
            </Link>
          </li>
        </ul>
      </div>
      <div className="flex items-center">
        <div className="px-4 py-2">
          <ModeToggle />
        </div>
        {props.auth ? (
          <DropdownMenu>
            <DropdownMenuTrigger>
              <Avatar
                avatarUrl={props.auth.gitlabAvatar?.avatarURL}
                fallbackUrl={props.auth.gitlabAvatar?.fallbackAvatarURL}
                name={props.auth.name!}
              />
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuLabel>User Account</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem>
                <UserIcon className="mr-2 h-4 w-4" />
                <span>Profile</span>
                <DropdownMenuShortcut>⇧P</DropdownMenuShortcut>
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Settings className="mr-2 h-4 w-4" />
                <span>Settings</span>
                <DropdownMenuShortcut>⇧S</DropdownMenuShortcut>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem>
                <LogOut className="mr-2 h-4 w-4" />
                <form method="POST" action="/api/v1/auth/sign-out">
                  <input type="hidden" name="csrf_token" value={csrfToken} />
                  <button type="submit" className="font-bold">
                    Log out
                  </button>
                </form>
                <DropdownMenuShortcut>⇧Q</DropdownMenuShortcut>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        ) : (
          <Link to="/login" search={{ redirect: location.href }}>
            Login
          </Link>
        )}
      </div>
    </nav>
  );
}
