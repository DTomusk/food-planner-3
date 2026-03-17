import { ChevronDownIcon, ChevronUpIcon } from "lucide-react";
import Inline from "../layout/Inline";
import IconButton from "../ui/IconButton";
import FormSectionTitle from "./FormSectionTitle";
import { useState } from "react";
import type { Status } from "../ui/StatusIcon";
import StatusIcon from "../ui/StatusIcon";

interface FormSectionProps {
    title: string;
    children: React.ReactNode;
    collapsible?: boolean;
    defaultCollapsed?: boolean;
    showStatus?: boolean;
    status?: Status;
}

export default function FormSection({ 
    title, 
    children, 
    collapsible = false,
    defaultCollapsed = false,
    showStatus = false,
    status,
}: FormSectionProps) {
    const [collapsed, setCollapsed] = useState(defaultCollapsed);

    return (
        <section className="space-y-4 bg-white border border-gray-200 rounded-lg p-6 shadow-sm">
            <div
                className={`
                ${collapsible ? "cursor-pointer select-none" : ""}
                `}
                onClick={() => collapsible && setCollapsed((prev) => !prev)}
            >
            <Inline justify="between" align="center">
                <FormSectionTitle title={title} />
                <Inline>
                    {showStatus && status && (
                        <StatusIcon status={status.status} />
                    )}
                {collapsible && (
                    <IconButton 
                        aria-label="Expand section" 
                        shape="circle" 
                        variant="primary-outline"
                        aria-expanded={!collapsed}>
                            {collapsed ? <ChevronDownIcon size={20} /> : <ChevronUpIcon size={20} />}
                    </IconButton>)}
                </Inline>
            </Inline>
        </div>
            {!collapsed && children}
        </section>
    );
}