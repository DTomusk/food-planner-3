import clsx from "clsx";
import { forwardRef, type TextareaHTMLAttributes } from "react";

export type TextAreaProps = TextareaHTMLAttributes<HTMLTextAreaElement>;

const TextArea = forwardRef<HTMLTextAreaElement, TextAreaProps>(
    ({ className, ...props }, ref) => {
        return(
            <div className="flex flex-col gap-1 w-full">
            <textarea
            ref={ref}
            className={clsx(
                "w-full rounded-md border px-3 py-2 text-sm",
                "focus:outline-none focus:ring-2 focus:ring-blue-500",
                "border-gray-300",
                className
            )}
            {...props}
            />
        </div>
        )
    });

TextArea.displayName = "TextArea";

export default TextArea;