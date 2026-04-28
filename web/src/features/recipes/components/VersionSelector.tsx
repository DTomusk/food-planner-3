import ActionList from "@/components/ui/ActionList";
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
};

export default function VersionSelector({ recipeId, versions, currentVersionNumber }: VersionSelectorProps) {
    const navigate = useNavigate();

    return (
        <aside className="w-full rounded-md border border-black bg-white p-3 md:w-56 md:shrink-0">
            <ActionList
                title="Recipe versions"
                items={versions.map((version) => {
                    const isCurrentVersion = version.version === currentVersionNumber;
                    return {
                        label: `Version ${version.version}${version.draft ? " (draft)" : ""}`,
                        onClick: () => navigate(`/recipes/${recipeId}/versions/${version.version}`),
                        selected: isCurrentVersion,
                    };
                })}
            />
        </aside>
    );
}