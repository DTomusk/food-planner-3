import Inline from "@/components/layout/Inline";
import { useMatches } from "react-router-dom";
import Breadcrumbs from "./Breadcrumbs";

type PageProps = {
    children: React.ReactNode;
    toolbarLeft?: React.ReactNode;
    toolbarActions?: React.ReactNode;
    toolbarClass?: string;
}

type CrumbHandle = {
    crumb?: (data?: unknown) => React.ReactNode;
}

export default function Page({children, toolbarLeft, toolbarActions, toolbarClass}: PageProps) {
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
                <div className="sticky top-16 z-30 border-b border-black bg-white/80 px-6 py-3 backdrop-blur sm:top-0">
                    <Inline justify="between" className="w-full">
                        <div className="flex items-center gap-4">
                        {hasBreadcrumbs && (
                            <Breadcrumbs crumbs={crumbs} />
                        )}

                        {toolbarLeft && <div>{toolbarLeft}</div>}
                        </div>

                        <Inline justify="end" align="center" className={toolbarClass}>
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