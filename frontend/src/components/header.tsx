export function Header({
  title,
  children,
}: {
  title: string;
  children?: React.ReactNode;
}) {
  return (
    <div className="flex flex-row justify-between">
      <h1 className="text-xl font-bold">{title}</h1>
      {children}
    </div>
  );
}
