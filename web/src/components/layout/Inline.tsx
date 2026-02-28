interface InlineProps {
    children: React.ReactNode;
    justify?: "start" | "center" | "end" | "between" | "around" | "evenly";
    align?: "start" | "center" | "end";
    wrap?: boolean;
}

export default function Inline({ children, justify = "center", align = "center", wrap = false }: InlineProps) {
    return <div className={`flex ${wrap ? "flex-wrap sm:flex-row flex-col" : "flex-nowrap flex-row"}  items-${align} justify-${justify} gap-2`}>{children}</div>;
}