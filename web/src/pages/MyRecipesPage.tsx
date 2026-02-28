import { Alert, PageTitle, Spinner } from "@/components";
import Container from "@/components/layout/Container";
import Stack from "@/components/layout/Stack";
import SectionTitle from "@/components/ui/SectionTitle";
import { RecipeList, useRecipes } from "@/features/recipes";
import { Page } from "@/layout";
import { extractErrorMessage } from "@/lib/errors";

export default function MyRecipesPage() {
    const { data, isLoading, error } = useRecipes();
    
    return (
        <Page>
            <PageTitle text="Recipes" />
            <Container size="md">
                <Stack space="lg">
            {isLoading && <Spinner/>}
            {error && <Alert message={extractErrorMessage(error)} closable />}
            {data && (
                <>
                    <SectionTitle text="My recipes" />
                    <RecipeList recipes={data.recipes} />
                </>
            )}
                </Stack>
            </Container>
        </Page>
    );
}