type HeaderProps = {
  title: string;
  size?: "xl" | "5xl";
  children?: React.ReactNode;
  margin?: string;
  subtitle?: string;
};

export function Header({ title, children, size = "xl", margin, subtitle }: HeaderProps) {
  console.log(margin);

  return (
    <div className="flex flex-row justify-between w-auto flex-wrap">
      <h1 className={`text-${size} font-bold ${margin} `}>{title}</h1>
      {subtitle && <h2 className="text-lg">{subtitle}</h2>}
      {children}
    </div>
  );
}
