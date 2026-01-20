interface ButtonProps {
    onClick?: (event: React.MouseEvent<HTMLButtonElement>) => void;
    disabled?: boolean;
    type?: "button" | "submit" | "reset";
    variant?: keyof typeof variants;
    children?: React.ReactNode;
    loading?: boolean;
}

const variants = {
    primary: "bg-blue-500 text-white rounded hover:bg-blue-600",
    secondary: "bg-gray-500 text-white rounded hover:bg-gray-600",
    danger: "bg-red-500 text-white rounded hover:bg-red-600",
}

export default function Button({ children, onClick, disabled, type = "button", variant = "primary", loading = false }: ButtonProps) {
    return (
    <button 
        onClick={onClick} 
        disabled={disabled || loading} 
        type={type} 
        className={`
            px-4 py-2
            ${variants[variant]}
            ${disabled || loading ? "opacity-50 cursor-not-allowed" : ""}
        `}>
            {loading ? "Loading..." : children}
    </button>
    );
}