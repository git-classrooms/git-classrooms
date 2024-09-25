import { cn } from "@/lib/utils";
import React from "react";

type HeaderProps = {
  title: React.ReactNode;
  children?: React.ReactNode;
  subtitle?: string;
} & Partial<Pick<HTMLDivElement, "className">>;

export function Header({ title, children, subtitle, className }: HeaderProps) {
  return (
    <div className="mb-10">
      <h1 className={cn("text-4xl tracking-tight font-bold mb-1", className)}>{title}</h1>
      {subtitle && <span className="text-muted-foreground">{subtitle}</span>}
      {children}
    </div>
  );
}
