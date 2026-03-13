import Inline from "@/components/layout/Inline";

type PageProps = {
    children: React.ReactNode;
    toolbarLeft?: React.ReactNode;
    toolbarActions?: React.ReactNode;
}

export default function Page({children, toolbarLeft, toolbarActions}: PageProps) {
    return (
        <>
            {(toolbarLeft || toolbarActions) && (
                <div className="sticky top-0 z-10 border-b border-gray-200 bg-white/90 backdrop-blur px-6 py-3">
                    <Inline justify="between" align="center" className="w-full">
                        <div>{toolbarLeft}</div>
                    <Inline justify="end" align="center">{toolbarActions}</Inline>
                    </Inline>
                </div>
            )}
            <div className="space-y-8 max-w-5xl mx-auto bg-white pt-12 px-6 pb-12">
                {children}
            </div>
        </>
    )
}