export default function Page({children} : {children: React.ReactNode}) {
    return (
        <div className="space-y-8 max-w-5xl mx-auto bg-white pt-8 px-6 pb-12 shadow">
            {children}
        </div>
    )
}