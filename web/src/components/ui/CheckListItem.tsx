import clsx from "clsx";

interface CheckListItemProps {
    checked: boolean;
    onChange: (checked: boolean) => void;
    children: React.ReactNode;
};

export default function CheckListItem({
    checked,
    onChange,
    children,
}: CheckListItemProps) {
    return (
        <label className="flex items-center gap-3 cursor-pointer select-none">
            <input
            type="checkbox"
            checked={checked}
            onChange={(e) => onChange(e.target.checked)}
            className="h-5 w-5 rounded border-gray-300 text-primary-500 focus:ring-primary-500"
            />

            <span
            className={clsx(
            "transition-all duration-200",
            checked && "line-through text-gray-400"
            )}
            >
            {children}
            </span>
        </label>
    );
}
