import clsx from "clsx";
import { forwardRef, type ButtonHTMLAttributes } from "react";

type IconButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
    variant?: "primary" | "secondary" | "danger" | "primary-outline" | "secondary-outline";
    shape?: "circle" | "square";
};

const IconButton = forwardRef<HTMLButtonElement, IconButtonProps>(
    ({ variant = "primary", shape = "square", className, type = "button", ...props }, ref) => {
        return (
            <button
                ref={ref}
                type={type}
                className={clsx(
                    "inline-flex items-center justify-center btn-icon",
                    variant === "primary" && "btn-primary",
                    variant === "secondary" && "btn-secondary",
                    variant === "danger" && "btn-danger",
                    variant === "primary-outline" && "btn-primary-outline",
                    variant === "secondary-outline" && "btn-secondary-outline",
                    shape === "circle" && "rounded-full",
                    shape === "square" && "rounded-md",
                    className
                )}
                {...props}
            />
        );
    }
);

IconButton.displayName = "IconButton";

export default IconButton;