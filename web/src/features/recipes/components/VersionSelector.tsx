import { Stack } from "@/components";
import ActionList from "@/components/ui/ActionList";
import Button from "@/components/ui/Button";
import CollapsibleHeader from "@/components/ui/CollapsibleHeader";
import { clsx } from "clsx";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

type Version = {
    version: number;
    createdAt: string;
    draft: boolean;
};

type VersionSelectorProps = {
    recipeId: string;
    versions: Version[];
    currentVersionNumber: number | undefined;
    canEdit: boolean;
    className?: string;
};

export default function VersionSelector({ recipeId, versions, currentVersionNumber, canEdit, className }: VersionSelectorProps) {
    const navigate = useNavigate();
    const hasDraft = versions.some((version) => version.draft);
    const [collapsed, setCollapsed] = useState(false);

    return (
        <aside className={clsx("w-full rounded-md border border-black bg-white p-3 md:w-56 md:shrink-0", className)}>
            <Stack space="md">
                <CollapsibleHeader collapsible title="Recipe versions" collapsed={collapsed} onToggle={() => setCollapsed(!collapsed)} />
                { collapsed ? null : (
                    <>
                        {canEdit ? (
                            <Button variant="secondary" onClick={() => navigate(`/recipes/${recipeId}/edit`)}>
                                {hasDraft ? "Edit draft" : "Create new version"}
                            </Button>
                        ) : null}
                    
                        <ActionList
                            items={versions.map((version) => {
                                const isCurrentVersion = version.version === currentVersionNumber;
                                return {
                                    label: `v${version.version}${version.draft ? " (draft)" : ""}`,
                                    onClick: () => navigate(`/recipes/${recipeId}/versions/${version.version}`),
                                    selected: isCurrentVersion,
                                };
                            })}
                        />
                    </>
                )}
            </Stack>
        </aside>
    );
}