import { useLocation, useNavigate, useParams } from "react-router-dom";
import { Alert, BackLink, PageTitle, Spinner } from "@/components";
import { useRecipe } from "@/features/recipes";
import { Page } from "@/layout";
import SharedBy from "@/components/SharedBy";
import IngredientList from "@/features/recipes/components/IngredientList";
import SectionTitle from "@/components/ui/SectionTitle";
import Container from "@/components/layout/Container";
import Stack from "@/components/layout/Stack";
import { useEffect, useState } from "react";
import { extractErrorMessage } from "@/lib/errors";
import MarkdownRenderer from "@/components/ui/MarkdownRenderer";
import { RecipeSourceType } from "@/features/recipes/types";
import IconButton from "@/components/ui/IconButton";
import { ClockFading, SquarePen } from "lucide-react";
import Dropdown from "@/components/ui/Dropdown";
import { useRecipeVersions } from "@/features/recipes/hooks/useRecipeVersions";

type RecipePageLocationState = {
    successMessage?: string;
};

export default function RecipePage() {
    const { id } = useParams<{ id: string }>();
    const recipeId = id ?? "";
    const location = useLocation();
    const navigate = useNavigate();
    const locationState = location.state as RecipePageLocationState | null;
    const [successMessage, setSuccessMessage] = useState<string | undefined>(locationState?.successMessage);

    const recipeQuery = useRecipe(recipeId);
    const recipe = recipeQuery.data?.recipe;
    const user = recipeQuery.data?.user;

    const versionsQuery = useRecipeVersions(recipeId);
    const versions = versionsQuery.data;

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
                            items: versions.map((version) => ({
                                label: `Version ${version.version} - ${new Date(version.createdAt).toLocaleString()}`,
                            })),
                        },
                    ]}
                />
            ) : null}
            <IconButton variant="primary-outline" onClick={() => navigate(`/recipes/${recipeId}/edit`)}>
                <SquarePen size={16} />
            </IconButton>
        </>
    );

    return (
        <Page toolbarLeft={<BackLink />} toolbarActions={toolbarActions}>
            <Container size="xl">
                <Stack space="xl">
                    {!recipeId && <Alert message="No recipe ID provided." />}
                    {recipeQuery.isLoading && <Spinner />}
                    {recipeQuery.error && <Alert message={extractErrorMessage(recipeQuery.error)} closable />}
                    {successMessage && <Alert message={successMessage} type="success" closable duration={3000} onClose={() => setSuccessMessage(undefined)} />}
                    {recipe ? (
                    <>
                    <PageTitle text={recipe.name} />
                    { user && <SharedBy user={user} /> }
                    <Stack space="sm">
                        <div className="text-center">Prep time: {recipe.prepMins} mins</div>
                        <div className="text-center">Cook time: {recipe.cookMins} mins</div>
                        <div className="text-center">Portions: {recipe.portions}</div>
                    </Stack>
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