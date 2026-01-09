export default function AppLayout({children}: {children: React.ReactNode}) {
    return (
        <div className="min-h-screen">
            {/* Add header etc. here */}
            <div className="flex flex-col bg-gray-100 min-h-screen">
                <main className="flex-1 p-6">
                    {children}
                </main>
            </div>
        </div>
    )
}