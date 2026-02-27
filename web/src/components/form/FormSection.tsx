import { ChevronDownIcon } from "lucide-react";
import Inline from "../layout/Inline";
import IconButton from "../ui/IconButton";
import FormSectionTitle from "./FormSectionTitle";

interface FormSectionProps {
    title: string;
    children: React.ReactNode;
    collapsible?: boolean;
}

export default function FormSection({ title, children, collapsible = false }: FormSectionProps) {
    return (
        <section className="space-y-4 bg-gray-50 border border-gray-200 rounded-lg p-6">
            <Inline justify="between">
                <FormSectionTitle title={title} />
                {collapsible && <IconButton aria-label="Expand section" shape="circle" variant="primary-outline"><ChevronDownIcon size={16} /></IconButton>}
            </Inline>
            {children}
        </section>
    );
}