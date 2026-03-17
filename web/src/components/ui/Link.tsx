interface LinkProps {
  onClick?: (e: React.MouseEvent<HTMLAnchorElement>) => void;
  color?: "black" | "primary";
  children: React.ReactNode;
}

export default function Link({ onClick, children, color = "black" }: LinkProps) {
  const colorClasses = {
    black: "text-black",
    primary: "text-primary-700",
  };

  return (
    <a
      onClick={onClick}
      className={`hover:underline cursor-pointer ${colorClasses[color]}`}
    >
      {children}
    </a>
  );
}