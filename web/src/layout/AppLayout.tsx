import NavBar from "./NavBar";

export default function AppLayout({children}: {children: React.ReactNode}) {
    return (
        <div className="min-h-screen bg-gray-100">
            <NavBar />
            <div className="flex flex-col min-h-screen">
                <main className="flex-1 p-6">
                    {children}
                </main>
            </div>
        </div>
    )
}