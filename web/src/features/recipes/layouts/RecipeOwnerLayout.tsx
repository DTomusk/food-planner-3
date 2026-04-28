import { Alert, Spinner } from "@/components";
import Container from "@/components/layout/Container";
import { useAuth } from "@/features/auth/hooks/useAuth";
import { useMe } from "@/features/users/hooks/useMe";
import { extractErrorMessage } from "@/lib/errors";
import { Navigate, Outlet, useParams } from "react-router-dom";
import { useRecipe } from "../hooks/useRecipe";

export default function RecipeOwnerLayout() {
    const { id } = useParams<{ id: string }>();
    const { isAuthenticated } = useAuth();
    const recipeId = id ?? "";
    const recipeQuery = useRecipe(recipeId);
    const meQuery = useMe();

    if (!isAuthenticated) {
        return <Navigate to="/auth/signin" replace />;
    }

    if (!recipeId) {
        return (
            <Container>
                <Alert message="No recipe ID provided." />
            </Container>
        );
    }

    if (recipeQuery.isLoading || meQuery.isLoading) {
        return (
            <Container>
                <Spinner />
            </Container>
        );
    }

    if (recipeQuery.error) {
        return (
            <Container>
                <Alert message={extractErrorMessage(recipeQuery.error)} closable />
            </Container>
        );
    }

    if (meQuery.error) {
        return (
            <Container>
                <Alert message={extractErrorMessage(meQuery.error)} closable />
            </Container>
        );
    }

    const recipeAuthorId = recipeQuery.data?.user.id;
    const currentUserId = meQuery.data?.id;

    if (!recipeAuthorId || !currentUserId || recipeAuthorId !== currentUserId) {
        return <Navigate to={`/recipes/${recipeId}`} replace />;
    }

    return <Outlet />;
}
