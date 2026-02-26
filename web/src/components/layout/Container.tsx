import clsx from "clsx";

interface ContainerProps {
    size?: "xs" | "sm" | "md" | "lg" | "xl" | "full";
    children: React.ReactNode;
}

export default function Container({ size = "md", children }: ContainerProps) {
    const sizes = {
        xs: "max-w-lg",
        sm: "max-w-xl",
        md: "max-w-2xl",
        lg: "max-w-4xl",
        xl: "max-w-6xl",
        full: "max-w-none",
    };
    return (
        <div className={clsx(`${sizes[size]} mx-auto px-4 sm:px-6 lg:px-8`)}>
            {children}
        </div>
    );
}