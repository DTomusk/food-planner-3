export default function Inline({ children }: { children: React.ReactNode }) {
    return <div className="flex items-center justify-center gap-2">{children}</div>;
}