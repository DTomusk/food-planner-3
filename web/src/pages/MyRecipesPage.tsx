import { Alert, PageTitle, Spinner } from "@/components";
import Container from "@/components/layout/Container";
import Stack from "@/components/layout/Stack";
import ConfirmModal from "@/components/ui/ConfirmModal";
import IconButton from "@/components/ui/IconButton";
import SectionTitle from "@/components/ui/SectionTitle";
import { RecipeList, useMyRecipes } from "@/features/recipes";
import { Page } from "@/layout";
import { extractErrorMessage } from "@/lib/errors";
import { Edit, Eye, Trash } from "lucide-react";
import { useState } from "react";

export default function MyRecipesPage() {
    const { data, isLoading, error } = useMyRecipes();
    const [selectedRecipeId, setSelectedRecipeId] = useState<string | null>(null);
    const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
    const [loading, setLoading] = useState(false);

    const handleDeleteClick = (recipeId: string) => {
        setSelectedRecipeId(recipeId);
        setShowDeleteConfirm(true);
    }

    const handleConfirmDelete = async () => {
        if (!selectedRecipeId) return;

        setLoading(true);
        try {
            console.log("Deleting recipe with ID:", selectedRecipeId);
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
            <PageTitle text="Recipes" />
            <Container size="md">
                <Stack space="lg">
            {isLoading && <Spinner/>}
            {error && <Alert message={extractErrorMessage(error)} closable />}
            {data && (
                <>
                    <SectionTitle text="My recipes" />
                    <RecipeList recipes={data.myRecipes}
                        onCardClick={()=>{}}
                        renderActions={(recipe) => (
                            <>
                            <IconButton variant="primary-outline" onClick={() => {}}>
                                <Eye size={16} />
                            </IconButton>
                            <IconButton variant="primary-outline" onClick={() => {}}>
                                <Edit size={16} />
                            </IconButton>
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
                />
            )}
        </Page>
    );
}