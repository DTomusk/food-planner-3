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
    primary: "bg-primary-600 text-white rounded hover:bg-primary-700",
    secondary: "bg-secondary-500 text-white rounded hover:bg-secondary-600",
    danger: "bg-red-500 text-white rounded hover:bg-red-600",
    primaryOutline: "border border-primary-600 text-primary-600 rounded hover:bg-primary-100 hover:text-primary-700",
    secondaryOutline: "border border-secondary-500 text-secondary-500 rounded hover:bg-secondary-100 hover:text-secondary-700",
    dangerOutline: "border border-red-500 text-red-500 rounded hover:bg-red-100 hover:text-red-700",
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