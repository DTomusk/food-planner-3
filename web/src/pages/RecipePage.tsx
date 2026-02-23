import { useParams } from "react-router-dom";
import { Alert, BackLink, PageTitle, Spinner } from "@/components";
import { useRecipe } from "@/features/recipes";
import { Page, Section } from "@/layout";
import SharedBy from "@/components/SharedBy";
import IngredientList from "@/features/recipes/components/IngredientList";
import SectionTitle from "@/components/ui/SectionTitle";

export default function RecipePage() {
    const { id } = useParams<{ id: string }>();
    const { data: recipe, isLoading, error } = useRecipe(id!);

    return (
        <Page>
            <BackLink to="/recipe" />
            {!id && <Alert message="No recipe ID provided." />}
            {isLoading && <Spinner />}
            {error && <Alert message={(error as Error).message} />}
            {recipe ? (
                <>
                <PageTitle text={recipe ? recipe.name : "Recipe Page"} />
                <SharedBy userName={recipe.user} />
                <Section>
                    <SectionTitle text="Ingredients" />
                    <IngredientList ingredients={recipe.ingredients} />
                </Section>
                </>
            ) : (
                !isLoading && id && !error && <Alert message="Recipe not found." />
            )}
        </Page>
    )
}