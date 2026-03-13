import { Alert, BackLink, PageTitle } from "@/components";
import Container from "@/components/layout/Container";
import { useIngredients } from "@/features/ingredients/hooks/useIngredients";
import { RecipeForm, useCreateRecipe, type RecipeFormValues } from "@/features/recipes";
import { mapFormValuesToCreateRecipeInput } from "@/features/recipes/mappers/recipeFormMapper";
import { Page } from "@/layout";
import { extractErrorMessage } from "@/lib/errors";
import { commonStrings } from "@/lib/strings";
import { useNavigate } from "react-router-dom";

export default function RecipeCreatePage() {
    const { mutate, isPending, error } = useCreateRecipe();
    const { data: ingredientsData } = useIngredients();
    const navigate = useNavigate();

    const handleSubmit = (values: RecipeFormValues) => {
      mutate(
          { input: mapFormValuesToCreateRecipeInput(values) },
          {
            onSuccess: (data) => {
              navigate(`/recipes/${data.createRecipe.id}`, {
                state: { successMessage: "Recipe created successfully!" }
              });
            }
          }
      );
    };

    return (
        <Page>
            <BackLink to="/recipes" />
            <PageTitle text={commonStrings.recipe.create} />
            {error && <Container><Alert message={extractErrorMessage(error)} closable /></Container>}
            <RecipeForm onSubmit={handleSubmit} isSubmitting={isPending} ingredients={ingredientsData || []} />
        </Page>
    );
}