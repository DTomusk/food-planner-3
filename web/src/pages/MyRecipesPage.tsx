import { Alert, PageTitle, Spinner } from "@/components";
import Container from "@/components/layout/Container";
import Stack from "@/components/layout/Stack";
import Text from "@/components/ui/Text";
import { RecipeList, useMyRecipes } from "@/features/recipes";
import { Page } from "@/layout";
import { extractErrorMessage } from "@/lib/errors";

export default function MyRecipesPage() {
    const { data, isLoading, error } = useMyRecipes();
    
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
                        renderActions={() => (
                            <>
                            {/* TODO: add these actions
                            <IconButton variant="primary-outline" onClick={() => {}}>
                                <Eye size={16} />
                            </IconButton>
                            <IconButton variant="primary-outline" onClick={() => {}}>
                                <Edit size={16} />
                            </IconButton> 
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