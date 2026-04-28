import { useNavigate, useParams } from "react-router-dom";
import { Alert, Button, Inline, Spinner } from "@/components";
import { Page } from "@/layout";
import Container from "@/components/layout/Container";
import Stack from "@/components/layout/Stack";
import { extractErrorMessage } from "@/lib/errors";
import IconButton from "@/components/ui/IconButton";
import { ClockFading } from "lucide-react";
import Dropdown from "@/components/ui/Dropdown";
import { useRecipeVersions } from "@/features/recipes/hooks/useRecipeVersions";
import { useRecipeVersion } from "@/features/recipes/hooks/useRecipeVersion";
import RecipeContentSections from "@/features/recipes/components/RecipeContentSections";
import VersionSelector from "@/features/recipes/components/VersionSelector";

export default function RecipeVersionPage() {
    const { id } = useParams<{ id: string }>();
    const { version } = useParams<{ version: string }>();
    const recipeId = id ?? "";
    const versionNumber = version ? parseInt(version, 10) : 0;
    const navigate = useNavigate();

    const recipeQuery = useRecipeVersion(recipeId, versionNumber);
    const recipe = recipeQuery.data?.recipe;
    const user = recipeQuery.data?.user;

    const versionsQuery = useRecipeVersions(recipeId);
    const versions = versionsQuery.data;

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
                            items: versions.map((version) => ({
                                label: `Version ${version.version} - ${new Date(version.createdAt).toLocaleString()}`,
                                onClick: () => navigate(`/recipes/${recipeId}/versions/${version.version}`),
                                disabled: version.version === versionNumber,
                            })),
                        },
                    ]}
                />
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
                    {recipe ? (
                        <Inline align="start" justify="center" className="w-full" gap="lg">
                            <RecipeContentSections recipe={recipe} user={user} />
                            {versions?.length ? (
                                <Stack>
                                    <VersionSelector
                                        recipeId={recipeId}
                                        versions={versions}
                                        currentVersionNumber={versionNumber}
                                    />
                                    <Button variant="secondary" onClick={() => navigate(`/recipes/${recipeId}/edit`)}>
                                        Create new version
                                    </Button>
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