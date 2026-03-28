import { Dialog, DialogPanel } from "@headlessui/react";
import Stack from "@/components/layout/Stack";
import IconButton from "@/components/ui/IconButton";
import { X } from "lucide-react";

interface MobileNavDrawerProps {
    open: boolean;
    onClose: () => void;
    children: React.ReactNode;
}

export default function MobileNavDrawer({ open, onClose, children }: MobileNavDrawerProps) {
    return (
        <Dialog open={open} onClose={onClose} className="relative z-50">
            <div className="fixed inset-0 bg-black/30" aria-hidden="true" />
            <DialogPanel className="fixed right-0 top-0 h-full w-64 bg-white px-4 py-6 shadow-lg">
                <Stack space="lg">
                    <div className="flex justify-end">
                        <IconButton onClick={onClose} variant="primary-outline" aria-label="Close menu">
                            <X size={16} />
                        </IconButton>
                    </div>
                    {children}
                </Stack>
            </DialogPanel>
        </Dialog>
    );
}
