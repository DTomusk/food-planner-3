export default function SubmitButton({ children, disabled }: { children: React.ReactNode; disabled?: boolean }) {
    return (
        <button
            type="submit"
            disabled={disabled}
            className={`px-4 py-2 bg-blue-500 text-white rounded ${disabled ? 'opacity-50 cursor-not-allowed' : 'hover:bg-blue-600'}`}
        >
            {children}
        </button>
    );
}