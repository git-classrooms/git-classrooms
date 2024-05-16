import { cn } from "@/lib/utils";

type HeaderProps = {
  title: string;
  children?: React.ReactNode;
  subtitle?: string;
} & Pick<HTMLDivElement, "className">;

export function Header({ title, children, subtitle, className }: HeaderProps) {
  return (
    <div>
      <h1 className={cn("text-xl font-bold mb-10", className)}>{title}</h1>
      {subtitle && <h2 className="text-lg">{subtitle}</h2>}
      {children}
    </div>
  );
}
