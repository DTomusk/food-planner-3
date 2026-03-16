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
                    "inline-flex items-center justify-center p-2 cursor-pointer",
                    variant === "primary" && "bg-primary-500 text-white hover:bg-primary-600",
                    variant === "secondary" && "bg-gray-500 text-white hover:bg-gray-600",
                    variant === "danger" && "bg-red-500 text-white hover:bg-red-600",
                    variant === "primary-outline" && "border border-primary-500 text-primary-500 hover:bg-primary-100 hover:text-primary-700",
                    variant === "secondary-outline" && "border border-gray-500 text-gray-500 hover:bg-gray-50",
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