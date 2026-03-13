import { Alert, PageTitle, Spinner } from "@/components";
import Container from "@/components/layout/Container";
import Stack from "@/components/layout/Stack";
import IconButton from "@/components/ui/IconButton";
import Text from "@/components/ui/Text";
import { RecipeList, useMyRecipes } from "@/features/recipes";
import { Page } from "@/layout";
import { extractErrorMessage } from "@/lib/errors";
import { Edit, Eye } from "lucide-react";
import { useNavigate } from "react-router-dom";

export default function MyRecipesPage() {
    const { data, isLoading, error } = useMyRecipes();
    const navigate = useNavigate();
    
    return (
        <Page>
            <PageTitle text="My recipes" />
            <Container size="md">
                <Stack space="lg">
            {isLoading && <Spinner/>}
            {error && <Alert message={extractErrorMessage(error)} closable />}
            {data && (
                <>
                    <Text>Here you can find all the recipes you've shared.</Text>
                    <RecipeList recipes={data}
                        onCardClick={()=>{}}
                        renderActions={(recipe) => (
                            <>
                            <IconButton variant="primary-outline" onClick={() => navigate(`/recipes/${recipe.id}`)}>
                                <Eye size={16} />
                            </IconButton>
                            <IconButton variant="primary-outline" onClick={() => navigate(`/recipes/${recipe.id}/edit`)}>
                                <Edit size={16} />
                            </IconButton> 
                            {/* TODO: add these actions
                            <IconButton variant="danger" onClick={() => handleDeleteClick(recipe.id)}>
                                <Trash size={16} />
                            </IconButton>*/}
                            </>
                        )} />
                </>
            )}
                </Stack>
            </Container>
        </Page>
    );
}