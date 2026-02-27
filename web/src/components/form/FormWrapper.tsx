export default function Form({children, onSubmit}: {children: React.ReactNode, onSubmit: (e: React.FormEvent<HTMLFormElement>) => void}) {
    return (
        <form onSubmit={onSubmit}>
            {children}
        </form>
    );
}