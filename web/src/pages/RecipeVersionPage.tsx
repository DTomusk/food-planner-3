import { useNavigate, useParams } from "react-router-dom";
import { Alert, Spinner } from "@/components";
import { Page } from "@/layout";
import IngredientList from "@/features/recipes/components/IngredientList";
import SectionTitle from "@/components/ui/SectionTitle";
import Container from "@/components/layout/Container";
import Stack from "@/components/layout/Stack";
import { extractErrorMessage } from "@/lib/errors";
import MarkdownRenderer from "@/components/ui/MarkdownRenderer";
import { RecipeSourceType } from "@/features/recipes/types";
import IconButton from "@/components/ui/IconButton";
import { ClockFading } from "lucide-react";
import Dropdown from "@/components/ui/Dropdown";
import { useRecipeVersions } from "@/features/recipes/hooks/useRecipeVersions";
import { useRecipeVersion } from "@/features/recipes/hooks/useRecipeVersion";
import RecipeDetailsCard from "@/features/recipes/components/RecipeDetailsCard";

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
            <Container size="md">
                <Stack space="xl">
                    {!recipeId && <Alert message="No recipe ID provided." />}
                    {recipeQuery.isLoading && <Spinner />}
                    {recipeQuery.error && <Alert message={extractErrorMessage(recipeQuery.error)} closable />}
                    {recipe ? (
                    <>
                    <RecipeDetailsCard
                        recipeTitle={recipe.name}
                        imageUrl={recipe.imageUrl}
                        description="blah blah blah"
                        prepTimeMinutes={recipe.prepMins}
                        cookTimeMinutes={recipe.cookMins}
                        portions={recipe.portions}
                        versionNumber={versionNumber}
                        sharedBy={ user ? { username: user.username, id: user.id } : undefined }
                    />
                    <Container size="xs">
                        <Stack space="lg">
                            <SectionTitle text="Ingredients" />
                            <IngredientList ingredients={recipe.ingredients} />
                            {recipe.source.type === RecipeSourceType.Website && recipe.source.url && (<>
                                <SectionTitle text="Website reference" />
                                <a href={recipe.source.url} target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline">
                                    {recipe.source.url}
                                </a>
                            </>)}
                            {recipe.source.type === RecipeSourceType.Cookbook && recipe.source.bookTitle && recipe.source.bookPage && (<>
                                <SectionTitle text="Book reference" />
                                <div>{recipe.source.bookTitle}, page {recipe.source.bookPage}</div>
                            </>)}
                            {recipe.source.type === RecipeSourceType.Original && recipe.source.instructions && (<>

                                <SectionTitle text="Instructions" />
                                <MarkdownRenderer content={recipe.source.instructions} />
                            </>)}
                        </Stack>
                    </Container>
                    </>
                ) : (
                    !recipeQuery.isLoading && recipeId && !recipeQuery.error && <Alert message="Recipe not found." />
                )}
                </Stack>
            </Container>
        </Page>
    );
}