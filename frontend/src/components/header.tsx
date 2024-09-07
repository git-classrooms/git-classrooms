import { cn } from "@/lib/utils";

type HeaderProps = {
  title: string;
  children?: React.ReactNode;
  subtitle?: string;
} & Partial<Pick<HTMLDivElement, "className">>;

export function Header({ title, children, subtitle, className }: HeaderProps) {
  return (
    <div className="mb-16">
      <h1 className={cn("text-4xl font-bold mb-1", className)}>{title}</h1>
      {subtitle && <span className="text-muted-foreground">{subtitle}</span>}
      {children}
    </div>
  );
}
