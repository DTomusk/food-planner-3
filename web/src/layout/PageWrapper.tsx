import Inline from "@/components/layout/Inline";
import { Link, useMatches } from "react-router-dom";

type PageProps = {
    children: React.ReactNode;
    toolbarLeft?: React.ReactNode;
    toolbarActions?: React.ReactNode;
}

type CrumbHandle = {
    crumb?: (data?: unknown) => React.ReactNode;
}

export default function Page({children, toolbarLeft, toolbarActions}: PageProps) {
    const matches = useMatches();

    const crumbs = matches
        .filter((match) => (match.handle as CrumbHandle)?.crumb)
        .map((match) => {
            const handle = match.handle as CrumbHandle;
            return {
                path: match.pathname,
                label: handle.crumb?.(match.loaderData),
            };
        });

    const hasBreadcrumbs = crumbs.length > 1;

    return (
        <>
            {(toolbarLeft || toolbarActions || hasBreadcrumbs) && (
                <div className="sticky top-0 z-10 border-b border-black bg-white/80 backdrop-blur px-6 py-3">
                    <Inline justify="between" className="w-full">
                        <div className="flex items-center gap-4">
                        {hasBreadcrumbs && (
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
                        )}

                        {toolbarLeft && <div>{toolbarLeft}</div>}
                        </div>

                        <Inline justify="end" align="center">
                        {toolbarActions}
                        </Inline>
                    </Inline>
                </div>
            )}
            <div className="space-y-8 max-w-5xl mx-auto pt-12 pb-12">
                {children}
            </div>
        </>
    )
}