export default function Section({ children }: { children: React.ReactNode }) {
    return (
        <div className="space-y-6 max-w-md mx-auto">
            {children}
        </div>
    )
}