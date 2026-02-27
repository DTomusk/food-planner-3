import { ChevronDownIcon, ChevronUpIcon } from "lucide-react";
import Inline from "../layout/Inline";
import IconButton from "../ui/IconButton";
import FormSectionTitle from "./FormSectionTitle";
import { useState } from "react";

interface FormSectionProps {
    title: string;
    children: React.ReactNode;
    collapsible?: boolean;
    defaultCollapsed?: boolean;
}

export default function FormSection({ 
    title, 
    children, 
    collapsible = false,
    defaultCollapsed = false 
}: FormSectionProps) {
    const [collapsed, setCollapsed] = useState(defaultCollapsed);

    return (
        <section className="space-y-4 bg-gray-50 border border-gray-200 rounded-lg p-6">
            <Inline justify="between">
                <FormSectionTitle title={title} />
                {collapsible && (
                    <IconButton 
                        aria-label="Expand section" 
                        shape="circle" 
                        variant="primary-outline"
                        onClick={() => setCollapsed((prev) => !prev)}
                        aria-expanded={!collapsed}>
                            {collapsed ? <ChevronDownIcon size={20} /> : <ChevronUpIcon size={20} />}
                    </IconButton>)}
            </Inline>
            {!collapsed && children}
        </section>
    );
}