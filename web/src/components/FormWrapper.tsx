export default function Form({children, onSubmit}: {children: React.ReactNode, onSubmit: (e: React.FormEvent<HTMLFormElement>) => void}) {
    return (
        <form onSubmit={onSubmit} className="space-y-4">
            {children}
        </form>
    );
}