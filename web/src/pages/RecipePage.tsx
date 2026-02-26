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

export default function RecipePage() {
    const { id } = useParams<{ id: string }>();
    const { data: recipe, isLoading, error } = useRecipe(id!);
    const location = useLocation();
    const navigate = useNavigate();
    const [successMessage, setSuccessMessage] = useState<string | undefined>(location.state?.successMessage);

    useEffect(() => {
        if (successMessage) {
            navigate(location.pathname, { replace: true, state: {} });
        }
    }, [successMessage, location.pathname, navigate]);

    return (
        <Page>
            <Container size="xl">
                <Stack space="xl">
                    <BackLink to="/recipe" />
                    <Container size="sm">
                        {!id && <Alert message="No recipe ID provided." />}
                        {isLoading && <Spinner />}
                        {error && <Alert message={extractErrorMessage(error)} closable />}
                        {successMessage && <Alert message={successMessage} type="success" closable duration={3000} onClose={() => setSuccessMessage(undefined)} />}
                    </Container>
                    {recipe ? (
                    <>
                    <PageTitle text={recipe ? recipe.name : "Recipe Page"} />
                    <SharedBy userName={recipe.user} />
                    <Stack space="sm">
                        <div className="text-center">Prep time: {recipe.prepMins} mins</div>
                        <div className="text-center">Cook time: {recipe.cookMins} mins</div>
                        <div className="text-center">Portions: {recipe.portions}</div>
                    </Stack>
                    <Container size="xs">
                        <Stack space="lg">
                            <SectionTitle text="Ingredients" />
                            <IngredientList ingredients={recipe.ingredients} />
                            {recipe.source.type === 1 && recipe.source.url && (<></>)}
                            {recipe.source.type === 2 && recipe.source.bookTitle && recipe.source.bookPage && (<></>)}
                            {recipe.source.instructions && (<>
                                <SectionTitle text="Instructions" />
                                <MarkdownRenderer content={recipe.source.instructions} />
                            </>)}
                        </Stack>
                    </Container>
                    </>
                ) : (
                    !isLoading && id && !error && <Alert message="Recipe not found." />
                )}
                </Stack>
            </Container>
        </Page>
    )
}