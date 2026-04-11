type FormSectionTitleProps = {
    title: string;
    size?: "sm" | "md" | "lg";
};


export default function FormSectionTitle({ title, size = "md" }: FormSectionTitleProps) {
    let className = "font-bold mb-4";
    switch (size) {
        case "sm":
            className += " text-sm";
            break;
        case "md":
            className += " text-md";
            break;
        case "lg":
            className += " text-lg";
            break;
    }
    return (
        <h2 className={className}>{title}</h2>
    );
}