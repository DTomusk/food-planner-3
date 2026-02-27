import { clsx } from "clsx";
import { Check, Clock, OctagonAlert } from "lucide-react";

export type Status = {
    status: "error" | "pending" | "completed";
};


export default function StatusIcon({ status }: Status) {
    return (
        <div
        className={clsx(
            "inline-flex items-center justify-center p-2 rounded-full border",
            status === "error" && "border-red-500 text-red-500",
            status === "pending" && "border-yellow-500 text-yellow-500",
            status === "completed" && "border-green-500 text-green-500",
        )}
        
        aria-label={status === "error" ? "Error" : status === "pending" ? "Pending" : "Completed"}>
            {status === "error" && <OctagonAlert size={20} />}
            {status === "pending" && <Clock size={20} />}
            {status === "completed" && <Check size={20} />}
        </div>
    );
}