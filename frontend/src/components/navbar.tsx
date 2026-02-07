import { ModeToggle } from "@/components/modeToggle";
import { Link, useRouterState } from "@tanstack/react-router";
import { useCsrf } from "@/provider/csrfProvider";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  ChevronDown,
  ExternalLink,
  GraduationCap,
  LogOut,
  Menu,
} from "lucide-react";
import { Avatar } from "./avatar";
import { GetMeResponse } from "@/swagger-client";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle, SheetTrigger } from "./ui/sheet";
import { cn } from "@/lib/utils";
import GitlabLogo from "@/assets/gitlab_logo.svg";

export function Navbar(props: { auth: GetMeResponse | null }) {
  return (
    <header className="sticky top-0 z-50 w-full border-b glass">
      <MobileNavbar {...props} />
      <DesktopNavbar {...props} />
    </header>
  );
}

interface NavbarProps {
  auth: GetMeResponse | null;
}

function DesktopNavbar(props: NavbarProps) {
  return (
    <nav className="hidden md:flex h-14 items-center justify-between px-6">
      <div className="flex items-center gap-6">
        <Logo />
        {props.auth && <NavLinks />}
      </div>
      <div className="flex items-center gap-2">
        <ModeToggle />
        {props.auth ? (
          <UserDropdown auth={props.auth} />
        ) : (
          <Button variant="glow" size="sm" asChild>
            <Link to="/login" search={{ redirect: location.href }}>
              Login
            </Link>
          </Button>
        )}
      </div>
    </nav>
  );
}

function Logo() {
  return (
    <Link to="/" className="flex items-center gap-2 group">
      <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-primary to-[hsl(280,100%,60%)] flex items-center justify-center transition-transform group-hover:scale-105">
        <GraduationCap className="w-5 h-5 text-primary-foreground" />
      </div>
      <span className="font-mono font-semibold text-lg tracking-tight hidden sm:block">
        GitClassrooms
      </span>
    </Link>
  );
}

function NavLinks() {
  const router = useRouterState();
  const currentPath = router.location.pathname;

  const links = [
    { to: "/dashboard", label: "Dashboard" },
    { to: "/classrooms", label: "Classrooms" },
  ] as const;

  return (
    <nav className="flex items-center gap-1">
      {links.map((link) => {
        const isActive = currentPath.startsWith(link.to);
        return (
          <Link
            key={link.to}
            to={link.to}
            className={cn(
              "relative px-3 py-2 text-sm font-medium transition-colors",
              "text-muted-foreground hover:text-foreground",
              isActive && "text-foreground"
            )}
          >
            {link.label}
            {isActive && (
              <span className="absolute bottom-0 left-3 right-3 h-0.5 bg-primary shadow-[0_0_8px_hsl(var(--glow-primary))]" />
            )}
          </Link>
        );
      })}
    </nav>
  );
}

function UserDropdown({ auth }: { auth: GetMeResponse }) {
  const { csrfToken } = useCsrf();

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <button
          className={cn(
            "flex items-center gap-2 pl-1 pr-2 py-1 rounded-full",
            "bg-muted/50 hover:bg-muted border border-border/50 hover:border-border",
            "transition-all duration-200",
            "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background"
          )}
        >
          <Avatar
            avatarUrl={auth.avatarURL}
            fallbackUrl={auth.fallbackAvatarURL}
            name={auth.name!}
            className="w-7 h-7"
          />
          <ChevronDown className="w-3.5 h-3.5 text-muted-foreground" />
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-64 p-0 overflow-hidden">
        {/* User Info Header */}
        <div className="p-4 bg-muted/30 border-b border-border/50">
          <div className="flex items-center gap-3">
            <Avatar
              avatarUrl={auth.avatarURL}
              fallbackUrl={auth.fallbackAvatarURL}
              name={auth.name!}
              className="w-10 h-10"
            />
            <div className="flex-1 min-w-0">
              <p className="font-medium truncate">{auth.name}</p>
              <p className="text-xs text-muted-foreground truncate">
                @{auth.gitlabUsername}
              </p>
            </div>
          </div>
        </div>

        {/* Menu Items */}
        <div className="p-1.5">
          <DropdownMenuItem asChild className="cursor-pointer rounded-md">
            <a
              href={auth.gitlabUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-3 px-3 py-2.5"
            >
              <div className="w-8 h-8 rounded-md bg-[#FC6D26]/10 flex items-center justify-center">
                <img src={GitlabLogo} className="w-4 h-4" alt="" />
              </div>
              <div className="flex-1">
                <p className="text-sm font-medium">GitLab Profile</p>
                <p className="text-xs text-muted-foreground">View your profile</p>
              </div>
              <ExternalLink className="w-3.5 h-3.5 text-muted-foreground" />
            </a>
          </DropdownMenuItem>
        </div>

        <DropdownMenuSeparator className="my-0" />

        {/* Logout */}
        <div className="p-1.5">
          <form method="POST" action="/api/v1/auth/sign-out">
            <input type="hidden" name="csrf_token" value={csrfToken} />
            <DropdownMenuItem asChild className="cursor-pointer rounded-md">
              <button
                type="submit"
                className="w-full flex items-center gap-3 px-3 py-2.5 text-destructive focus:text-destructive"
              >
                <div className="w-8 h-8 rounded-md bg-destructive/10 flex items-center justify-center">
                  <LogOut className="w-4 h-4" />
                </div>
                <div className="flex-1 text-left">
                  <p className="text-sm font-medium">Log out</p>
                  <p className="text-xs opacity-70">End your session</p>
                </div>
              </button>
            </DropdownMenuItem>
          </form>
        </div>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

function MobileNavbar(props: NavbarProps) {
  return (
    <div className="md:hidden flex h-14 items-center justify-between px-4">
      <Logo />
      <div className="flex items-center gap-2">
        <ModeToggle />
        {props.auth ? (
          <>
            <UserDropdown auth={props.auth} />
            <Sheet>
              <SheetTrigger asChild>
                <Button variant="ghost" size="icon">
                  <Menu className="h-5 w-5" />
                  <span className="sr-only">Open menu</span>
                </Button>
              </SheetTrigger>
              <SheetContent side="right">
                <SheetHeader>
                  <SheetTitle className="flex items-center gap-2">
                    <Logo />
                  </SheetTitle>
                  <SheetDescription className="sr-only">Navigation menu</SheetDescription>
                </SheetHeader>
                <nav className="flex flex-col gap-2 mt-6">
                  <Link
                    to="/dashboard"
                    className="flex items-center px-3 py-2 text-sm font-medium rounded-md hover:bg-muted transition-colors"
                  >
                    Dashboard
                  </Link>
                  <Link
                    to="/classrooms"
                    className="flex items-center px-3 py-2 text-sm font-medium rounded-md hover:bg-muted transition-colors"
                  >
                    Classrooms
                  </Link>
                </nav>
              </SheetContent>
            </Sheet>
          </>
        ) : (
          <Button variant="glow" size="sm" asChild>
            <Link to="/login" search={{ redirect: location.href }}>
              Login
            </Link>
          </Button>
        )}
      </div>
    </div>
  );
}
