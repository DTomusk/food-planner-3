import Link from "./Link";

type NavItemProps = {
    onClick?: () => void;
    label: string;
    icon?: React.ReactNode;
}

export default function NavItem({ onClick, label, icon }: NavItemProps) {
    return (
        <Link onClick={onClick}>{icon} {label}</Link>
    );
}