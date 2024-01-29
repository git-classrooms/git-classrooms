import React from "react";

export function TopBar(): JSX.Element {
  return (
    <nav className="top-bar">
      <div className="top-bar__logo">
        {/* <Link to="/">GitLab</Link> */}</div>
      <div className="top-bar__links">
        {/* 
        <Link to="/projects">Projects</Link>
        <Link to="/groups">Groups</Link>
        <Link to="/issues">Issues</Link>
        <Link to="/merge-requests">Merge Requests</Link>
        */}
      </div>
      <div className="top-bar__user">
        <span className="top-bar__user-avatar">Your Avatar</span>
        <span className="top-bar__user-name">Your Name</span>
      </div>
    </nav>
  );
}
