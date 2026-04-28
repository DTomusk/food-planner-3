import { useLocation, useNavigate, useParams } from "react-router-dom";
import { Alert, Button, Inline, Spinner } from "@/components";
import { useAuth } from "@/features/auth/hooks/useAuth";
import { useRecipe } from "@/features/recipes";
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
    const { id } = useParams<{ id: string }>();
    const recipeId = id ?? "";
    const location = useLocation();
    const navigate = useNavigate();
    const locationState = location.state as RecipePageLocationState | null;
    const [successMessage, setSuccessMessage] = useState<string | undefined>(locationState?.successMessage);

    const recipeQuery = useRecipe(recipeId);
    const recipe = recipeQuery.data?.recipe;
    const user = recipeQuery.data?.user;
    const meQuery = useMe({ enabled: isAuthenticated });
    const canEditRecipe = Boolean(user?.id && meQuery.data?.id && user.id === meQuery.data.id);

    const versionsQuery = useRecipeVersions(recipeId);
    const versions = versionsQuery.data;
    const currentVersionNumber = versions?.length
        ? Math.max(...versions.map((version) => version.version))
        : undefined;

    useEffect(() => {
        if (!successMessage) {
            return;
        }

        navigate(location.pathname, { replace: true, state: {} });
    }, [successMessage, location.pathname, navigate]);

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
                                const isCurrentVersion = version.version === currentVersionNumber;

                                return {
                                    label: `Version ${version.version} - ${new Date(version.createdAt).toLocaleString()}${isCurrentVersion ? " (current)" : ""}`,
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
        <Page toolbarActions={toolbarActions}>
            <Container size="xl">
                <Stack space="xl">
                    {!recipeId && <Alert message="No recipe ID provided." />}
                    {recipeQuery.isLoading && <Spinner />}
                    {recipeQuery.error && <Alert message={extractErrorMessage(recipeQuery.error)} closable />}
                    {successMessage && <Alert message={successMessage} type="success" closable duration={3000} onClose={() => setSuccessMessage(undefined)} />}
                    {recipe ? (
                        <Inline align="start" justify="center" className="w-full" gap="lg">
                            <RecipeContentSections recipe={recipe} user={user} />
                            {versions?.length ? (
                                <Stack>
                                    <VersionSelector
                                        recipeId={recipeId}
                                        versions={versions}
                                        currentVersionNumber={currentVersionNumber}
                                    />
                                    {canEditRecipe ? (
                                        <Button variant="secondary" onClick={() => navigate(`/recipes/${recipeId}/edit`)}>
                                            Create new version
                                        </Button>
                                    ) : null}
                                </Stack>
                            ) : null}
                        </Inline>
                ) : (
                    !recipeQuery.isLoading && recipeId && !recipeQuery.error && <Alert message="Recipe not found." />
                )}
                </Stack>
            </Container>
        </Page>
    );
}