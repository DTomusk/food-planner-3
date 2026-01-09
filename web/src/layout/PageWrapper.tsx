export default function Page({children} : {children: React.ReactNode}) {
    return (
        <div className="space-y-8 max-w-7xl mx-auto bg-white">
            {children}
        </div>
    )
}