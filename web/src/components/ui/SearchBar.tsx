import { Search } from "lucide-react";

type SearchBarProps = {
    value: string;
    onChange: (value: string) => void;
    placeholder?: string;
    onSubmit?: () => void;
    loading?: boolean;
    disabled?: boolean;
};

export default function SearchBar(
    { value, onChange, placeholder = "Search...", onSubmit, loading = false, disabled = false }: SearchBarProps
) {
    const isDisabled = disabled || loading;

    return (
        <form
            className="relative"
            onSubmit={(e) => {
                e.preventDefault();
                if (onSubmit) {
                    onSubmit();
                }
            }}
        >
            <input
                type="text"
                value={value}
                onChange={(e) => onChange(e.target.value)}
                placeholder={placeholder}
                aria-label="Search"
                disabled={isDisabled}
                className="w-full rounded-md border bg-white border-black px-4 py-2 pr-10 focus-ring disabled:cursor-not-allowed disabled:bg-gray-100"
            />

            <button
                type="submit"
                aria-label="Submit search"
                disabled={isDisabled}
                className="absolute inset-y-0 right-0 inline-flex w-10 items-center justify-center text-gray-500 transition hover:text-gray-700 focus-ring disabled:cursor-not-allowed disabled:text-gray-300"
            >
                <Search size={16} />
            </button>
        </form>
    );
}