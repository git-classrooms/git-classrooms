import React from "react";
import { Header } from "./header";

interface LayoutProps {
  title: string;
  subtitle?: string;
  children: React.ReactNode;
}

export const Layout: React.FC<LayoutProps> = ({ children, title, subtitle }) => (
  <div className="mx-6 md:px-10">
    <Header title={title} subtitle={subtitle} className="text-5xl" />
    <div className="grid grid-cols-1 lg:grid-cols-2 justify-between gap-10">{children}</div>
  </div>
);
