interface TextProps {
  children: React.ReactNode;
  variant?: "body" | "muted" | "caption" | "error";
  as?: React.ElementType;
  align?: "left" | "center" | "right";
}

export default function Text({
  children,
  variant = "body",
  as: Component = "p",
  align = "left",
}: TextProps) {
  const variantStyles = {
    body: "text-base text-gray-900",
    muted: "text-sm text-gray-500",
    caption: "text-xs text-gray-400",
    error: "text-sm text-red-600",
  };

  return (
    <Component className={`${variantStyles[variant]} text-${align}`}>
      {children}
    </Component>
  );
}