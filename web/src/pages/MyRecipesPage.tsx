import { Alert, PageTitle, Spinner } from "@/components";
import Container from "@/components/layout/Container";
import Stack from "@/components/layout/Stack";
import ConfirmModal from "@/components/ui/ConfirmModal";
import IconButton from "@/components/ui/IconButton";
import Link from "@/components/ui/Link";
import Text from "@/components/ui/Text";
import { RecipeList, useMyRecipes } from "@/features/recipes";
import { useDeleteRecipe } from "@/features/recipes/hooks/useDeleteRecipe";
import { Page } from "@/layout";
import { extractErrorMessage } from "@/lib/errors";
import { Trash } from "lucide-react";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

export default function MyRecipesPage() {
    const { data, isLoading, error } = useMyRecipes();
    const { mutate: deleteRecipe, error: deleteError } = useDeleteRecipe();
    const [selectedRecipeId, setSelectedRecipeId] = useState<string | null>(null);
    const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    const handleDeleteClick = (recipeId: string) => {
        setSelectedRecipeId(recipeId);
        setShowDeleteConfirm(true);
    }

    const handleConfirmDelete = async () => {
        if (!selectedRecipeId) return;

        setLoading(true);
        try {
            deleteRecipe({ input: { id: selectedRecipeId } });
            setShowDeleteConfirm(false);
        } catch (err) {
            console.error("Error deleting recipe:", err);
        } finally {
            setLoading(false);
            setSelectedRecipeId(null);
        }
    };
    
    return (
        <Page>
            <PageTitle text="My recipes" />
            <Container size="md">
                <Stack space="lg">
            {isLoading && <Spinner/>}
            {error && <Alert message={extractErrorMessage(error)} closable />}
            {data && !data.me && <Alert message="You need to be logged in to view your recipes." closable />}
            {data && data.me && (
                <>
                    <Text>Here you can find all the recipes you've shared. Have you deleted a recipe by mistake or otherwise want to bring one back from the dead? Then click <Link text="here" color="primary" onClick={()=> navigate("deleted")}/>.</Text>
                    <RecipeList recipes={data.me.recipes}
                        onCardClick={()=>{}}
                        renderActions={(recipe) => (
                            <>
                            {/* TODO: add these actions
                            <IconButton variant="primary-outline" onClick={() => {}}>
                                <Eye size={16} />
                            </IconButton>
                            <IconButton variant="primary-outline" onClick={() => {}}>
                                <Edit size={16} />
                            </IconButton> */}
                            <IconButton variant="danger" onClick={() => handleDeleteClick(recipe.id)}>
                                <Trash size={16} />
                            </IconButton>
                            </>
                        )} />
                </>
            )}
                </Stack>
            </Container>
            {showDeleteConfirm && (
                <ConfirmModal
                    isOpen={showDeleteConfirm}
                    title="Confirm Delete"
                    description="Are you sure you want to delete this recipe? Note: deleted recipes can be restored within 30 days of deletion, after that, they're gone forever."
                    confirmText="Delete"
                    cancelText="Cancel"
                    loading={loading}
                    onConfirm={handleConfirmDelete}
                    onCancel={() => setShowDeleteConfirm(false)}
                    variant="danger"
                    error={deleteError ? extractErrorMessage(deleteError) : undefined}
                />
            )}
        </Page>
    );
}