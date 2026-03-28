import Inline from "@/components/layout/Inline";
import clsx from "clsx";

type NavItemProps = {
    onClick?: () => void;
    label: string;
    icon?: React.ReactNode;
}

export default function NavItem({ onClick, label, icon }: NavItemProps) {
    return (
        <button
            type="button"
            onClick={onClick}
            className={clsx(
                "w-full rounded-md border border-black bg-primary-50 px-3 py-3 text-left text-sm font-medium text-primary-900 transition-colors duration-150",
                "hover:bg-primary-100 focus-ring",
                "cursor-pointer hover:scale-105 transition-transform",
            )}
        >
            <Inline justify="start" align="center" gap="sm" className="w-full">
                {icon && <span aria-hidden="true" className="shrink-0 text-primary-900 [&>svg]:size-4">{icon}</span>}
                <span>{label}</span>
            </Inline>
        </button>
    );
}