import { Alert, Button, PageTitle, Spinner } from "@/components";
import { useAuth } from "@/features/auth/hooks/useAuth";
import { RecipeList, useRecipes } from "@/features/recipes";
import { Page } from "@/layout";
import { useNavigate } from "react-router-dom";

export default function RecipeListingPage() {
    const {data, isLoading, error} = useRecipes();
    const { isAuthenticated } = useAuth();
    
    const navigate = useNavigate();
    return (
        <Page>
            <PageTitle text="Recipes" />
            {!isAuthenticated && <Alert type="info" message="Please sign in to add a new recipe." />}
            {isAuthenticated && <Button onClick={() => navigate("/recipe/create")}>
                Add New Recipe
            </Button>}
            {isLoading && <Spinner/>}
            {error && <Alert message={(error as Error).message} />}
            {data && (
            <RecipeList recipes={data.recipes} />
            )}
        </Page>
    );
}