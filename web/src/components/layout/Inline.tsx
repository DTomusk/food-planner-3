interface InlineProps {
    children: React.ReactNode;
    justify?: "start" | "center" | "end" | "between" | "around" | "evenly";
    align?: "start" | "center" | "end";
}

export default function Inline({ children, justify = "center", align = "center" }: InlineProps) {
    return <div className={`flex items-${align} justify-${justify} gap-2`}>{children}</div>;
}