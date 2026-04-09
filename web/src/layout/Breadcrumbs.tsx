import type { ReactNode } from "react";
import { Link } from "react-router-dom";

export default function Breadcrumbs({ crumbs }: { crumbs: { path: string; label: ReactNode }[] }) {
    return (
        <nav className="text-sm text-gray-600">
        {crumbs.map((crumb, index) => {
            const isLast = index === crumbs.length - 1;

        return (
        <span key={crumb.path}>
            {!isLast ? (
            <>
                <Link to={crumb.path}>{crumb.label}</Link> /{" "}
            </>
            ) : (
            <span className="text-gray-900 font-medium">
                {crumb.label}
            </span>
            )}
        </span>
        );
    })}
    </nav>
    );
}