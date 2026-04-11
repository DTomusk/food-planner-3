import { X } from "lucide-react";

export default function FilterChip({ label, onClear }: { label: string; onClear: () => void }) {
    return (
        <button
            className="inline-flex items-center gap-1 rounded bg-primary-100 border border-primary-200 px-2 py-1 text-sm text-gray-800 hover:bg-primary-200 hover:cursor-pointer"
            onClick={onClear}
            aria-label={`Remove filter: ${label}`}
        >
            {label} <span aria-hidden="true"><X size={12} /></span>
        </button>
    );
}