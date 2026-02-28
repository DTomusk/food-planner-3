interface LinkProps {
    onClick: () => void;
    text: string;
    color?: "black" | "primary";
}

export default function Link({ onClick, text, color = "black" }: LinkProps) {
    const colorClasses = {
        black: "text-black",
        primary: "text-blue-600",
    };

    return <a 
    onClick={onClick}
    className={`hover:underline cursor-pointer ${colorClasses[color]}`}
    >{text}</a>;
}