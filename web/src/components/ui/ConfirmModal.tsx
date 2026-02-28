import { Dialog, DialogPanel, DialogTitle } from "@headlessui/react";
import { Alert, Button } from "@/components";

interface ConfirmModalProps {
    isOpen: boolean;
    title: string;
    description?: string;
    confirmText?: string;
    cancelText?: string;
    loading?: boolean;
    onConfirm: () => Promise<void> | void;
    onCancel: () => void;
    variant?: "primary" | "danger";
    error?: string;
}

export default function ConfirmModal({
    isOpen,
    title,
    description,
    confirmText = "Confirm",
    cancelText = "Cancel",
    loading = false,
    onConfirm,
    onCancel,
    variant = "primary",
    error
}: ConfirmModalProps) {
    return (
        <Dialog open={isOpen} onClose={onCancel} className="relative z-50">
            <div className="fixed inset-0 flex w-screen items-center justify-center bg-black/50">
                <DialogPanel className="max-w-lg space-y-4 border bg-white p-6 rounded">
                    <DialogTitle className="text-2xl font-bold">{title}</DialogTitle>
                    {error && <Alert message={error} closable />}
                    {description && <p>{description}</p>}
                    <div className="flex justify-end gap-4">
                        <Button onClick={onCancel} disabled={loading} variant="primaryOutline">{cancelText}</Button>
                        <Button onClick={onConfirm} disabled={loading} variant={variant}>{confirmText}</Button>
                    </div>
                </DialogPanel>
            </div>
        </Dialog>
    );
}