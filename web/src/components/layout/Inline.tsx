import clsx from "clsx";

interface InlineProps {
    children: React.ReactNode;
    justify?: "start" | "center" | "end" | "between" | "around" | "evenly";
    align?: "start" | "center" | "end" | "none";
    wrap?: boolean;
    gap?: "none" | "sm" | "md" | "lg";
    className?: string;
}

export default function Inline({
    children,
    justify = "center",
    align = "none",
    wrap = false,
    gap = "sm",
    className,
}: InlineProps) {
    const justifyClasses = {
        start: "justify-start",
        center: "justify-center",
        end: "justify-end",
        between: "justify-between",
        around: "justify-around",
        evenly: "justify-evenly",
    };

    const alignClasses = {
        start: "items-start",
        center: "items-center",
        end: "items-end",
        none: "",
    };

    const gapClasses = {
        none: "gap-0",
        sm: "gap-2",
        md: "gap-4",
        lg: "gap-6",
    };

    return (
        <div
            className={clsx(
                "flex",
                wrap ? "flex-wrap sm:flex-row flex-col" : "flex-nowrap flex-row",
                justifyClasses[justify],
                alignClasses[align],
                gapClasses[gap],
                className,
            )}
        >
            {children}
        </div>
    );
}