import { Stack, Text } from "@/components";
import ActionList from "@/components/ui/ActionList";
import Button from "@/components/ui/Button";
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
    hasDraft: boolean;
};

export default function VersionSelector({ recipeId, versions, currentVersionNumber, canEdit, hasDraft }: VersionSelectorProps) {
    const navigate = useNavigate();

    return (
        <aside className="w-full rounded-md border border-black bg-white p-3 md:w-56 md:shrink-0">
            <Stack space="md">
                <Text>Recipe versions</Text>
                {canEdit ? (
                    <Button variant="secondary" onClick={() => navigate(`/recipes/${recipeId}/edit`)}>
                        {hasDraft ? "Edit draft" : "Create new version"}
                    </Button>
                ) : null}
                <ActionList
                    items={versions.map((version) => {
                        const isCurrentVersion = version.version === currentVersionNumber;
                        return {
                            label: `Version ${version.version}${version.draft ? " (draft)" : ""}`,
                            onClick: () => navigate(`/recipes/${recipeId}/versions/${version.version}`),
                            selected: isCurrentVersion,
                        };
                    })}
                />
            </Stack>
        </aside>
    );
}