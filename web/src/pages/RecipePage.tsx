import { useLocation, useNavigate, useParams } from "react-router-dom";
import { Alert, Inline, Spinner } from "@/components";
import { useAuth } from "@/features/auth/hooks/useAuth";
import { useRecipe } from "@/features/recipes";
import { isRecipeVersionUnavailableError, useRecipeVersion } from "@/features/recipes/hooks/useRecipeVersion";
import { useMe } from "@/features/users/hooks/useMe";
import { Page } from "@/layout";
import Container from "@/components/layout/Container";
import Stack from "@/components/layout/Stack";
import { useEffect, useState } from "react";
import { extractErrorMessage } from "@/lib/errors";
import IconButton from "@/components/ui/IconButton";
import { ClockFading, SquarePen } from "lucide-react";
import Dropdown from "@/components/ui/Dropdown";
import { useRecipeVersions } from "@/features/recipes/hooks/useRecipeVersions";
import RecipeContentSections from "@/features/recipes/components/RecipeContentSections";
import VersionSelector from "@/features/recipes/components/VersionSelector";

type RecipePageLocationState = {
    successMessage?: string;
};

export default function RecipePage() {
    const { isAuthenticated } = useAuth();
    const { id, version } = useParams<{ id: string; version?: string }>();
    const recipeId = id ?? "";
    const selectedVersionParam = version ? Number.parseInt(version, 10) : undefined;
    const isViewingSpecificVersion =
        selectedVersionParam !== undefined && !Number.isNaN(selectedVersionParam);
    const location = useLocation();
    const navigate = useNavigate();
    const locationState = location.state as RecipePageLocationState | null;
    const [successMessage, setSuccessMessage] = useState<string | undefined>(locationState?.successMessage);

    const recipeQuery = useRecipe(recipeId, { enabled: !isViewingSpecificVersion });
    const recipeVersionQuery = useRecipeVersion(recipeId, selectedVersionParam ?? 0, {
        enabled: isViewingSpecificVersion,
    });
    const activeRecipeQuery = isViewingSpecificVersion ? recipeVersionQuery : recipeQuery;
    const recipe = activeRecipeQuery.data?.recipe;
    const user = activeRecipeQuery.data?.user;
    const meQuery = useMe({ enabled: isAuthenticated });
    const canEditRecipe = Boolean(user?.id && meQuery.data?.id && user.id === meQuery.data.id);

    const versionsQuery = useRecipeVersions(recipeId);
    const versions = versionsQuery.data;
    const selectedVersionNumber = isViewingSpecificVersion
        ? selectedVersionParam
        : recipe?.version;

    useEffect(() => {
        if (!successMessage) {
            return;
        }

        navigate(location.pathname, { replace: true, state: {} });
    }, [successMessage, location.pathname, navigate]);

    useEffect(() => {
        if (!isViewingSpecificVersion) {
            return;
        }

        if (!recipeVersionQuery.error || !isRecipeVersionUnavailableError(recipeVersionQuery.error)) {
            return;
        }

        navigate(`/recipes/${recipeId}`, { replace: true });
    }, [isViewingSpecificVersion, navigate, recipeId, recipeVersionQuery.error]);

    const toolbarActions = (
        <>
            {versions?.length ? (
                <Dropdown
                    button={
                        <IconButton variant="primary-outline">
                            <ClockFading size={16} />
                        </IconButton>
                    }
                    sections={[
                        {
                            title: "Recipe versions",
                            items: versions.map((version) => {
                                const isCurrentVersion = version.version === selectedVersionNumber;

                                return {
                                    label: `v${version.version}${version.draft ? " (draft)" : ""}`,
                                    onClick: () => navigate(`/recipes/${recipeId}/versions/${version.version}`),
                                    disabled: isCurrentVersion,
                                };
                            }),
                        },
                    ]}
                />
            ) : null}
            {canEditRecipe ? (
                <IconButton variant="primary-outline" onClick={() => navigate(`/recipes/${recipeId}/edit`)}>
                    <SquarePen size={16} />
                </IconButton>
            ) : null}
        </>
    );

    return (
        <Page toolbarActions={toolbarActions} toolbarClass="block lg:hidden">
            <Container size="xl">
                <Stack space="xl">
                    {!recipeId && <Alert message="No recipe ID provided." />}
                    {activeRecipeQuery.isLoading && <Spinner />}
                    {activeRecipeQuery.error && !isRecipeVersionUnavailableError(activeRecipeQuery.error) && <Alert message={extractErrorMessage(activeRecipeQuery.error)} closable />}
                    {successMessage && <Alert message={successMessage} type="success" closable duration={3000} onClose={() => setSuccessMessage(undefined)} />}
                    {recipe ? (
                        <Inline align="start" justify="center" className="w-full" gap="lg">
                            <RecipeContentSections recipe={recipe} user={user} />
                            {versions?.length ? (
                                <VersionSelector
                                    recipeId={recipeId}
                                    versions={versions}
                                    currentVersionNumber={selectedVersionNumber}
                                    canEdit={canEditRecipe}
                                    className="hidden lg:block"
                                />
                            ) : null}
                        </Inline>
                ) : (
                    !activeRecipeQuery.isLoading && recipeId && !activeRecipeQuery.error && <Alert message="Recipe not found." />
                )}
                </Stack>
            </Container>
        </Page>
    );
}