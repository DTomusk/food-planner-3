import type { ReactNode } from "react";

export default function Column({ children }: { children: ReactNode }) {
    return (
        <div className="flex flex-col gap-4 items-center">
            {children}
        </div>
    );
}