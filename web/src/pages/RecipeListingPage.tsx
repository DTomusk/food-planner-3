import { Alert, PageTitle, Spinner } from "@/components";
import Container from "@/components/layout/Container";
import Inline from "@/components/layout/Inline";
import Stack from "@/components/layout/Stack";
import IconButton from "@/components/ui/IconButton";
import SectionTitle from "@/components/ui/SectionTitle";
import { useAuth } from "@/features/auth/hooks/useAuth";
import { RecipeList, useRecipes } from "@/features/recipes";
import { Page } from "@/layout";
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
                    <IconButton onClick={() => navigate("/recipe/create")} aria-label="Add new recipe">
                        <Plus size={40}/>
                    </IconButton>
                    <p className="text-sm">Add Recipe</p>
                </Stack>
            </Inline>}
            {isLoading && <Spinner/>}
            {error && <Alert message={(error as Error).message} />}
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