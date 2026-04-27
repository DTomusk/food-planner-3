export default function Form({children, onSubmit}: {children: React.ReactNode, onSubmit?: (e: React.FormEvent<HTMLFormElement>) => void}) {
    return (
        <form onSubmit={(e) => {
            e.preventDefault();
            onSubmit?.(e);
        }}>
            {children}
        </form>
    );
}