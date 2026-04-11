import { ChevronDownIcon, ChevronUpIcon } from "lucide-react";
import FormSectionTitle from "../form/FormSectionTitle";
import Inline from "../layout/Inline";
import IconButton from "./IconButton";
import StatusIcon, { type Status } from "./StatusIcon";

interface CollapsibleHeaderProps {
    title: string;
    collapsible?: boolean;
    defaultCollapsed?: boolean;
    showStatus?: boolean;
    status?: Status;
    titleSize?: "sm" | "md" | "lg";
    onToggle?: (collapsed: boolean) => void;
    collapsed?: boolean; // controlled collapsed state
}

export default function CollapsibleHeader({
    title,
    collapsible = false,
    showStatus = false,
    status,
    titleSize = "md",
    onToggle,
    collapsed,
}: CollapsibleHeaderProps) {
    return (
        <Inline justify="between" align="center">
            <FormSectionTitle title={title} size={titleSize} />
            <Inline>
                {showStatus && status && (
                    <StatusIcon status={status.status} />
                )}
            {collapsible && (
                <IconButton 
                    aria-label={collapsed ? "Expand section" : "Collapse section"}
                    shape="circle" 
                    variant="primary-outline"
                    aria-expanded={!collapsed}
                    onClick={() => onToggle?.(!collapsed)}>

                        {collapsed ? <ChevronDownIcon size={20} /> : <ChevronUpIcon size={20} />}
                </IconButton>)}
            </Inline>
        </Inline>
    )
}