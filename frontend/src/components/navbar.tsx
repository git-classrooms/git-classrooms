import GitlabLogo from "@/assets/gitlab_logo.svg";
import ReactLogo from "@/assets/react.svg";
import { ModeToggle } from "@/components/modeToggle.tsx";
import { Link } from "@tanstack/react-router";

export function Navbar() {
  return (
    <nav className="flex justify-between px-8 py-2.5">
      <div className="flex items-center">
        <a href="/" className="">
          <img className="h-14" src={GitlabLogo} alt="Gitlab Logo" />
        </a>
        <ul className="flex">
          <li className="content-center">
            <Link to="/" className="font-medium text-sm px-4 py-2 hover:underline"
                  activeProps={{ className: "!font-bold" }}>Dashboard</Link>
          </li>
          <li className="content-center">
            <Link to="/" className="font-medium text-sm px-4 py-2 hover:underline"
                  activeProps={{ className: "!font-bold" }}>Created Classrooms</Link>
          </li>
          <li className="content-center">
            <Link to="/" className="font-medium text-sm px-4 py-2 hover:underline"
                  activeProps={{ className: "!font-bold" }}>Joined Classrooms</Link>
          </li>
        </ul>
      </div>
      <div className="flex items-center">
        <div className="px-4 py-2">
          <ModeToggle />
        </div>
        <img className="h-10 mr-2" src={ReactLogo} alt="User Image" />
      </div>
    </nav>
  );
}
