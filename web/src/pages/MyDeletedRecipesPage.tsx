import { Alert, BackLink, PageTitle, Spinner } from "@/components";
import Container from "@/components/layout/Container";
import Stack from "@/components/layout/Stack";
import Text from "@/components/ui/Text";
import { RecipeList, useMyDeletedRecipes } from "@/features/recipes";
import { Page } from "@/layout";
import { extractErrorMessage } from "@/lib/errors";

export default function MyDeletedRecipesPage() {
    const { data, isLoading, error } = useMyDeletedRecipes();
    
    return (
        <Page>
            <BackLink to="/recipes/me" />
            <PageTitle text="Deleted recipes" />
            <Container size="md">
                <Stack space="lg">
                <Text>View any of your deleted recipes below. You can choose to restore any of the listed recipes. Recipes that have been deleted for 30 days or more will be permanently deleted.</Text>
                {isLoading && <Spinner/>}
                {error && <Alert message={extractErrorMessage(error)} closable />}
                {data && (
                    <RecipeList recipes={data.myDeletedRecipes} onCardClick={()=>{}} />
                )}
                </Stack>
            </Container>
        </Page>
    );
}