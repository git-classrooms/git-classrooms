import { cn } from "@/lib/utils";

type HeaderProps = {
  title: string;
  children?: React.ReactNode;
  subtitle?: string;
} & Partial<Pick<HTMLDivElement, "className">>;

export function Header({ title, children, subtitle, className }: HeaderProps) {
  return (
    <div className="mb-10">
      <h1 className={cn("text-xl font-bold mb-3", className)}>{title}</h1>
      {subtitle && <h2 className="text-gray-400 dark:text-gray-500">{subtitle}</h2>}
      {children}
    </div>
  );
}
