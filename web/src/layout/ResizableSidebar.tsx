import { useEffect, useRef, useState, type CSSProperties, type PointerEvent, type ReactNode } from "react";

type ResizableSidebarProps = {
    children: ReactNode;
    minWidth?: number;
    maxWidth?: number;
    defaultWidth?: number;
    className?: string;
};

const DEFAULT_MIN_WIDTH = 192;
const DEFAULT_MAX_WIDTH = 384;
const DEFAULT_WIDTH = 256;

// TODO: move to components if used elsewhere
export default function ResizableSidebar({
    children,
    minWidth = DEFAULT_MIN_WIDTH,
    maxWidth = DEFAULT_MAX_WIDTH,
    defaultWidth = DEFAULT_WIDTH,
    className = "",
}: ResizableSidebarProps) {
    const initialWidth = Math.min(maxWidth, Math.max(minWidth, defaultWidth));
    const [width, setWidth] = useState(initialWidth);
    const [isResizing, setIsResizing] = useState(false);
    const resizeStartX = useRef(0);
    const resizeStartWidth = useRef(initialWidth);

    useEffect(() => {
        if (!isResizing) {
            return;
        }

        const previousUserSelect = document.body.style.userSelect;
        const previousCursor = document.body.style.cursor;

        document.body.style.userSelect = "none";
        document.body.style.cursor = "col-resize";

        const onPointerMove = (event: globalThis.PointerEvent) => {
            const deltaX = event.clientX - resizeStartX.current;
            const nextWidth = Math.min(
                maxWidth,
                Math.max(minWidth, resizeStartWidth.current + deltaX)
            );

            setWidth(nextWidth);
        };

        const onPointerUp = () => {
            setIsResizing(false);
        };

        window.addEventListener("pointermove", onPointerMove);
        window.addEventListener("pointerup", onPointerUp);

        return () => {
            window.removeEventListener("pointermove", onPointerMove);
            window.removeEventListener("pointerup", onPointerUp);
            document.body.style.userSelect = previousUserSelect;
            document.body.style.cursor = previousCursor;
        };
    }, [isResizing, maxWidth, minWidth]);

    const handleResizeStart = (event: PointerEvent<HTMLDivElement>) => {
        if (event.button !== 0) {
            return;
        }

        resizeStartX.current = event.clientX;
        resizeStartWidth.current = width;
        setIsResizing(true);
        event.preventDefault();
    };

    return (
        <aside
            className={`relative hidden sm:block sm:[width:var(--sidebar-width)] ${className}`}
            style={{ "--sidebar-width": `${width}px` } as CSSProperties}
        >
            {children}
            <div
                className="absolute right-0 top-0 h-full w-1 translate-x-1/2 cursor-col-resize bg-transparent transition-colors hover:bg-black/20"
                onPointerDown={handleResizeStart}
                aria-hidden
            />
        </aside>
    );
}