import { Alert, BackLink, PageTitle, Spinner } from "@/components";
import Container from "@/components/layout/Container";
import { useIngredients } from "@/features/ingredients/hooks/useIngredients";
import { RecipeForm, useRecipe, type RecipeFormValues } from "@/features/recipes";
import { useUpdateRecipe } from "@/features/recipes/hooks/useUpdateRecipe";
import { mapFormValuesToCreateRecipeInput } from "@/features/recipes/mappers/recipeFormMapper";
import { Page } from "@/layout";
import { commonStrings } from "@/lib";
import { extractErrorMessage } from "@/lib/errors";
import { useNavigate, useParams } from "react-router-dom";

export default function RecipeUpdatePage() {
    const { id } = useParams<{ id: string }>();
    const { data, isLoading, error } = useRecipe(id ?? "");
    const { mutate, isPending, error: mutateError } = useUpdateRecipe();
    const { data: ingredientsData } = useIngredients();
    const navigate = useNavigate();

    const handleSubmit = (values: RecipeFormValues) => {
        if (!id) {
            return;
        }

        mutate(
            {
                input: {
                    id,
                    details: mapFormValuesToCreateRecipeInput(values),
                },
            },
            {
            onSuccess: (data) => {
                navigate(`/recipes/${data.updateRecipe.id}`, {
                state: { successMessage: "Recipe updated successfully!" }
                });
            }
            }
        );
    };

    return (
        <Page>
            <BackLink />
            <PageTitle text={commonStrings.recipe.update} />
            {!id && <Container><Alert message="No recipe ID provided." /></Container>}
            {isLoading && <Container><Spinner /></Container>}
            {error && <Container><Alert message={extractErrorMessage(error)} closable /></Container>}
            {mutateError && <Container><Alert message={extractErrorMessage(error)} closable /></Container>}
            {id && data?.formValues && !isLoading && (
                <RecipeForm
                    key={id}
                    onSubmit={handleSubmit}
                    isSubmitting={isPending}
                    ingredients={ingredientsData || []}
                    defaultValues={data.formValues}
                />
            )}
        </Page>
    );
}