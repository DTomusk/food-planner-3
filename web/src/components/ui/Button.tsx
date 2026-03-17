import clsx from "clsx";

interface ButtonProps {
  onClick?: (event: React.MouseEvent<HTMLButtonElement>) => void;
  disabled?: boolean;
  type?: "button" | "submit" | "reset";
  variant?: 
    | "primary"
    | "secondary"
    | "danger"
    | "primaryOutline"
    | "secondaryOutline"
    | "dangerOutline";
  children?: React.ReactNode;
  loading?: boolean;
}

export default function Button({
  children,
  onClick,
  disabled,
  type = "button",
  variant = "primary",
  loading = false,
}: ButtonProps) {
  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled || loading}
      className={clsx(
        "btn",
        {
          "btn-primary": variant === "primary",
          "btn-secondary": variant === "secondary",
          "btn-danger": variant === "danger",
          "btn-primary-outline": variant === "primaryOutline",
          "btn-secondary-outline": variant === "secondaryOutline",
          "btn-danger-outline": variant === "dangerOutline",
          "btn-disabled": disabled || loading,
        }
      )}
    >
      {loading ? "Loading..." : children}
    </button>
  );
}