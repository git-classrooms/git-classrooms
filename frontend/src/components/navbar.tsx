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
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { LogOut, User as UserIcon } from "lucide-react";
import { Avatar } from "./avatar";
import { GetMeResponse } from "@/swagger-client";
import { Button } from "@/components/ui/button.tsx";


export function Navbar(props: { auth: GetMeResponse | null }) {
  const { csrfToken } = useCsrf();

  return (
    <nav className="flex justify-between px-8 py-2.5 mb-8 border-b">
      <div className="flex items-center">
        <a href="/" className="">
          <img className="h-14" src={GitlabLogo} alt="Gitlab Logo" />
        </a>
        {props.auth ? (
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
        ) : (<div />)}
      </div>
      <div className="flex items-center">
        <div className="px-4 py-2">
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
              <DropdownMenuLabel>
                <div className="font-medium">{props.auth.name}</div>
                <div className="text-sm text-muted-foreground md:inline">
                  @{props.auth.gitlabUsername}
                </div>
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem className="flex items-center" onKeyDown={(event) => {
                if (event.key === " " || event.key === "Spacebar") {
                  event.preventDefault();
                  window.open(props.auth?.gitlabUrl, "_blank");
                }
              }}>
                <a href={props.auth.gitlabUrl} target="_blank" rel="noopener noreferrer" className="flex items-center">
                  <UserIcon className="mr-2 h-4 w-4" />
                  <span>Profile</span>
                </a>
              </DropdownMenuItem>
              <DropdownMenuItem onSelect={event => event.preventDefault()}>
                <ModeToggle/>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem onKeyDown={(event) => {
                if (event.key === " " || event.key === "Spacebar") {
                  (document.getElementById("logOutForm") as HTMLFormElement).submit();
                }
              }}>
                <LogOut className="mr-2 h-4 w-4" />
                <form id="logOutForm" method="POST" action="/api/v1/auth/sign-out">
                  <input type="hidden" name="csrf_token" value={csrfToken} />
                  <button type="submit" className="font-bold">
                    Log out
                  </button>
                </form>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        ) : (
          <div className="flex">
            <ModeToggle />
            <Button className="ml-3">
              <Link to="/login" search={{ redirect: location.href }}>
                Login
              </Link>
            </Button>
          </div>
        )}
      </div>
    </nav>
  );
}
