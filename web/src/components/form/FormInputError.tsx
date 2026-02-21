import { AlertTriangle } from "lucide-react";

export default function FormInputError({ message }: { message?: string }) {
    if (!message) {
        return <p className="min-h-[1.25rem]"></p>;
    }
    return <p className="text-sm font-medium text-red-500 flex items-center gap-1"><AlertTriangle size={16} />{message}</p>;
}