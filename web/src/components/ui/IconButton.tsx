import clsx from "clsx";
import { forwardRef, type ButtonHTMLAttributes } from "react";

type IconButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
    variant?: "primary" | "secondary" | "danger";
};

const IconButton = forwardRef<HTMLButtonElement, IconButtonProps>(
    ({ variant = "primary", className, type = "button", ...props }, ref) => {
        return (
            <button
                ref={ref}
                type={type}
                className={clsx(
                    "inline-flex items-center justify-center rounded-md p-2",
                    variant === "primary" && "bg-blue-500 text-white",
                    variant === "secondary" && "bg-gray-500 text-white",
                    variant === "danger" && "bg-red-500 text-white",
                    className
                )}
                {...props}
            />
        );
    }
);

IconButton.displayName = "IconButton";

export default IconButton;