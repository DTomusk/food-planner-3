import clsx from "clsx";
import { forwardRef, type SelectHTMLAttributes } from "react";
import Text from "./Text";

export type SelectProps = SelectHTMLAttributes<HTMLSelectElement> & {
    error?: string;
};

const Select = forwardRef<HTMLSelectElement, SelectProps>(
    ({ className, error, children, ...props }, ref) => {
        return (
            <div className="flex flex-col gap-1">
                <select
                    ref={ref}
                    className={clsx(
                        "w-full rounded-md border px-3 py-2 text-sm bg-white",
                        "focus-ring",
                        error
                        ? "border-red-500 focus-ring-error"
                        : "border-gray-300",
                        className,
                        props.disabled && "cursor-not-allowed opacity-50"
                    )}
                    {...props}
                >
                    {children}
                </select>
                {error && <Text variant="error">{error}</Text>}
            </div>
        );
    }
);

Select.displayName = "Select";

export default Select;