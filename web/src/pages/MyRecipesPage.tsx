import { Alert, PageTitle, Spinner, Button, Inline } from "@/components";
import Stack from "@/components/layout/Stack";
import IconButton from "@/components/ui/IconButton";
import Text from "@/components/ui/Text";
import { RecipeList, useMyRecipes } from "@/features/recipes";
import { Page } from "@/layout";
import { extractErrorMessage } from "@/lib/errors";
import { Edit, Eye, Plus } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { commonStrings } from "@/lib/strings";

export default function MyRecipesPage() {
    const { data, isLoading, error } = useMyRecipes();
    const navigate = useNavigate();
    
    return (
        <Page>
            <PageTitle text="My recipes" />
            <div className="max-w-md md:max-w-2xl mx-auto px-4 sm:px-6 lg:px-8">
            <Stack space="lg">
            {isLoading && <Spinner/>}
            {error && <Alert message={extractErrorMessage(error)} closable />}
            {data && (
                <>
                    <Inline>
                        <Stack space="sm">
                            <Button onClick={() => navigate("/recipes/create")} 
                            aria-label="Add new recipe" variant="primary">
                                <Inline>
                                <Plus /> {commonStrings.recipe.add}
                                </Inline>
                            </Button>
                        </Stack>
                    </Inline>
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
            </div>
        </Page>
    );
}