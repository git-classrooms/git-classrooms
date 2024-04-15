import GitlabLogo from "@/assets/gitlab_logo.svg";
import ReactLogo from "@/assets/react.svg";
import { ModeToggle } from "@/components/modeToggle.tsx";

export function Navbar() {
  return(
    <nav className="flex justify-between px-8 py-2.5">
      <div className="flex items-center">
        <a  href="/" className="">
          <img className="h-14" src={GitlabLogo} alt="Gitlab Logo" />
        </a>
        <ul className="flex">
          <li className="content-center">
            <a href="/" className="font-medium text-sm px-4 py-2 hover:underline">Dashboard</a>
          </li>
          <li className="content-center">
            <a href="/" className="font-medium text-sm px-4 py-2 hover:underline">Created Classrooms</a>
          </li>
          <li className="content-center">
            <a href="/" className="font-medium text-sm px-4 py-2 hover:underline">Joined Classrooms</a>
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
