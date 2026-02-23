import clsx from "clsx";

interface StackProps {
  space?: "sm" | "md" | "lg" | "xl";
  children: React.ReactNode;
};

export default function Stack({ space = "md", children }: StackProps) {
    const spaces = {
        sm: "space-y-2",
        md: "space-y-4",
        lg: "space-y-8",
        xl: "space-y-12",
    };

    return <div className={clsx(spaces[space])}>{children}</div>;
}
