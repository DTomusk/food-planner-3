export default function Tag({ children }: { children: React.ReactNode }) {
    return (
        <span className="inline-flex w-fit rounded-full bg-primary-100 px-2 py-1 text-xs font-semibold text-primary-800">
            {children}
         </span>
    );
}