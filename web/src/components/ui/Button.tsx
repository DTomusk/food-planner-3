import clsx from "clsx";

interface ButtonProps {
    onClick?: (event: React.MouseEvent<HTMLButtonElement>) => void;
    disabled?: boolean;
    type?: "button" | "submit" | "reset";
    variant?: keyof typeof variants;
    children?: React.ReactNode;
    loading?: boolean;
}

const variants = {
    primary: "bg-primary-500 text-white rounded hover:bg-primary-600",
    secondary: "bg-gray-500 text-white rounded hover:bg-gray-600",
    danger: "bg-red-500 text-white rounded hover:bg-red-600",
    primaryOutline: "border border-primary-500 text-primary-500 rounded hover:bg-primary-100 hover:text-primary-700",
    secondaryOutline: "border border-gray-500 text-gray-500 rounded hover:bg-gray-50",
    dangerOutline: "border border-red-500 text-red-500 rounded hover:bg-red-50",
}

export default function Button({ children, onClick, disabled, type = "button", variant = "primary", loading = false }: ButtonProps) {
    return (
    <button 
        onClick={onClick} 
        disabled={disabled || loading} 
        type={type} 
        className={clsx(`
            px-4 py-2 text-sm font-medium cursor-pointer
            ${variants[variant]}
            ${disabled || loading ? "opacity-50 cursor-not-allowed" : ""}
        `)}>
            {loading ? "Loading..." : children}
    </button>
    );
}