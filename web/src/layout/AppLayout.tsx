import Inline from "@/components/layout/Inline";
import NavBar from "./NavBar";

export default function AppLayout({children}: {children: React.ReactNode}) {
    return (
        <div className="min-h-screen w-full bg-primary">
            <Inline wrap justify="start" align="start" gap="none" className="min-h-screen w-full">
                <NavBar />
                <main className="min-h-screen flex-1 min-w-0 overflow-x-hidden">
                    {children}
                </main>
            </Inline>
        </div>
    )
}