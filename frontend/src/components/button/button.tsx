interface ButtonProps {
  onClick: () => void;
  text: string;
}

export function Button({ onClick, text }: ButtonProps): JSX.Element {
  return <button onClick={onClick}>{text}</button>;
}
