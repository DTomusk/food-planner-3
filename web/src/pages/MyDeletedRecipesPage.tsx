import { Alert, BackLink, PageTitle, Spinner } from "@/components";
import Container from "@/components/layout/Container";
import Stack from "@/components/layout/Stack";
import ConfirmModal from "@/components/ui/ConfirmModal";
import IconButton from "@/components/ui/IconButton";
import Text from "@/components/ui/Text";
import { RecipeList, useMyDeletedRecipes } from "@/features/recipes";
import { useUndeleteRecipe } from "@/features/recipes/hooks/useUndeleteRecipe";
import { Page } from "@/layout";
import { extractErrorMessage } from "@/lib/errors";
import { HeartPlus } from "lucide-react";
import { useState } from "react";

export default function MyDeletedRecipesPage() {
    const { data, isLoading, error } = useMyDeletedRecipes();
    const [selectedRecipeId, setSelectedRecipeId] = useState<string | null>(null);
    const [showUndeleteConfirm, setShowUndeleteConfirm] = useState(false);
    const [loading, setLoading] = useState(false);
    const { mutate: undeleteRecipe, error: undeleteError } = useUndeleteRecipe();

    const handleUndeleteClick = (recipeId: string) => {
        setSelectedRecipeId(recipeId);
        setShowUndeleteConfirm(true);
    }

    const handleConfirmUndelete = async () => {
        if (!selectedRecipeId) return;

        setLoading(true);
        try {
            undeleteRecipe({ input: { id: selectedRecipeId } });
            setShowUndeleteConfirm(false);
        } catch (err) {
            console.error("Error undeleting recipe:", err);
        } finally {
            setLoading(false);
            setSelectedRecipeId(null);
        }
    };

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
                    <RecipeList recipes={data.myDeletedRecipes} 
                        onCardClick={()=>{}} 
                        renderActions={(recipe) => (
                           <IconButton variant="primary" onClick={() => handleUndeleteClick(recipe.id)}>
                                <HeartPlus size={16} />
                            </IconButton> 
                        )}/>
                )}
                </Stack>
            </Container>
            {showUndeleteConfirm && (
                <ConfirmModal
                    isOpen={showUndeleteConfirm}
                    title="Confirm Undelete"
                    description="Are you sure you want to restore this recipe? We promise you can delete it again later if you like 👉👈."
                    confirmText="Restore"
                    cancelText="Cancel"
                    loading={loading}
                    onConfirm={handleConfirmUndelete}
                    onCancel={() => setShowUndeleteConfirm(false)}
                    variant="primary"
                    error={undeleteError ? extractErrorMessage(undeleteError) : undefined}
                />
            )}
        </Page>
    );
}