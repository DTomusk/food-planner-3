import Inline from "@/components/layout/Inline";
import NavBar from "./NavBar";

export default function AppLayout({children}: {children: React.ReactNode}) {
    return (
        <div className="min-h-screen w-full bg-linear-to-br from-background to-white">
            <Inline wrap justify="start" align="none" gap="none" className="min-h-screen w-full">
                <NavBar />
                <main className="flex-1 min-w-0">
                    {children}
                </main>
            </Inline>
        </div>
    )
}