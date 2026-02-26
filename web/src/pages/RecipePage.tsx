import { useParams } from "react-router-dom";
import { Alert, BackLink, PageTitle, Spinner } from "@/components";
import { useRecipe } from "@/features/recipes";
import { Page } from "@/layout";
import SharedBy from "@/components/SharedBy";
import IngredientList from "@/features/recipes/components/IngredientList";
import SectionTitle from "@/components/ui/SectionTitle";
import Container from "@/components/layout/Container";
import Stack from "@/components/layout/Stack";

export default function RecipePage() {
    const { id } = useParams<{ id: string }>();
    const { data: recipe, isLoading, error } = useRecipe(id!);

    return (
        <Page>
            <Container size="xl">
                <Stack space="xl">
                <BackLink to="/recipe" />
                {!id && <Alert message="No recipe ID provided." />}
                {isLoading && <Spinner />}
                {error && <Alert message={(error as Error).message} />}
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