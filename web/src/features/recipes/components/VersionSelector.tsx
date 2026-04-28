import clsx from "clsx";
import { useNavigate } from "react-router-dom";

type Version = {
    version: number;
    createdAt: string;
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
            <p className="px-1 py-2 text-xs font-semibold text-gray-500 uppercase select-none">Recipe versions</p>
            <ul>
                {versions.map((version) => {
                    const isCurrentVersion = version.version === currentVersionNumber;
                    return (
                        <li key={version.version}>
                            <button
                                type="button"
                                disabled={isCurrentVersion}
                                onClick={() => navigate(`/recipes/${recipeId}/versions/${version.version}`)}
                                className={clsx(
                                    "flex w-full items-center gap-2 rounded px-2 py-2 text-left text-sm text-gray-900",
                                    "hover:bg-gray-100",
                                    isCurrentVersion && "cursor-not-allowed opacity-50",
                                    !isCurrentVersion && "cursor-pointer"
                                )}
                            >
                                {`Version ${version.version}${isCurrentVersion ? " (latest)" : ""}`}
                            </button>
                        </li>
                    );
                })}
            </ul>
        </aside>
    );
}