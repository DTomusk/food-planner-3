import { Alert, Button, PageTitle, Spinner } from "@/components";
import Container from "@/components/layout/Container";
import Inline from "@/components/layout/Inline";
import Stack from "@/components/layout/Stack";
import SectionTitle from "@/components/ui/SectionTitle";
import { useAuth } from "@/features/auth/hooks/useAuth";
import { RecipeList, useRecipes } from "@/features/recipes";
import { Page } from "@/layout";
import { extractErrorMessage } from "@/lib/errors";
import { commonStrings } from "@/lib/strings";
import { Plus } from "lucide-react";
import { useNavigate } from "react-router-dom";

export default function RecipeListingPage() {
    const { data, isLoading, error } = useRecipes();
    const { isAuthenticated } = useAuth();
    
    const navigate = useNavigate();
    return (
        <Page>
            <PageTitle text="Recipes" />
            <Container size="md">
                <Stack space="lg">
            {!isAuthenticated && <Alert type="info" message="Please sign in to add a new recipe." />}
            {isAuthenticated && 
            <Inline>
                <Stack space="sm">
                    <Button onClick={() => navigate("/recipes/create")} 
                    aria-label="Add new recipe" variant="primary">
                        <Inline>
                        <Plus /> {commonStrings.recipe.add}
                        </Inline>
                    </Button>
                </Stack>
            </Inline>}
            {isLoading && <Spinner/>}
            {error && <Alert message={extractErrorMessage(error)} closable />}
            {data && (
                <>
                    <SectionTitle text="Top recipes" />
                    <RecipeList recipes={data.recipes} />
                </>
            )}
                </Stack>
            </Container>
        </Page>
    );
}