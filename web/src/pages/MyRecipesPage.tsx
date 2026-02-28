import { Alert, PageTitle, Spinner } from "@/components";
import Container from "@/components/layout/Container";
import Stack from "@/components/layout/Stack";
import IconButton from "@/components/ui/IconButton";
import SectionTitle from "@/components/ui/SectionTitle";
import { RecipeList, useMyRecipes } from "@/features/recipes";
import { Page } from "@/layout";
import { extractErrorMessage } from "@/lib/errors";
import { Edit, Eye, Trash, View } from "lucide-react";

export default function MyRecipesPage() {
    const { data, isLoading, error } = useMyRecipes();
    
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
                            <IconButton variant="danger" onClick={() => {}}>
                                <Trash size={16} />
                            </IconButton>
                            </>
                        )} />
                </>
            )}
                </Stack>
            </Container>
        </Page>
    );
}