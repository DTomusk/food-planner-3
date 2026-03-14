import clsx from "clsx";

interface StackProps {
  space?: "sm" | "md" | "lg" | "xl";
  children: React.ReactNode;
  className?: string;
};

export default function Stack({ space = "md", children, className }: StackProps) {
    const spaces = {
        sm: "space-y-2",
        md: "space-y-4",
        lg: "space-y-8",
        xl: "space-y-12",
    };

    return <div className={clsx("flex flex-col", spaces[space], className)}>{children}</div>;
}
