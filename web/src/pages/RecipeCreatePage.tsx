import { Alert, BackLink, PageTitle } from "@/components";
import Container from "@/components/layout/Container";
import { useIngredients } from "@/features/ingredients/hooks/useIngredients";
import { RecipeForm, useCreateRecipe, type RecipeFormValues } from "@/features/recipes";
import { Page } from "@/layout";
import { extractErrorMessage } from "@/lib/errors";
import { commonStrings } from "@/lib/strings";
import { useNavigate } from "react-router-dom";

export default function RecipeCreatePage() {
    const { mutate, isPending, error } = useCreateRecipe();
    const { data: ingredientsData } = useIngredients();
    const navigate = useNavigate();

    const handleSubmit = (values: RecipeFormValues) => {
      console.log("Submitting recipe with values:", values);
      mutate(
          { input: { 
            name: values.name, 
            ingredientUsages: values.ingredientUsages.map(usage => ({
              ingredientID: usage.ingredientId,
              quantity: usage.quantity,
              unit: usage.unit
            })),
            prepMins: values.prepMins,
            cookMins: values.cookMins,
            portions: values.portions,
            recipeSource: {
              type: values.sourceType,
              url: values.sourceType === 1 ? values.url : undefined,
              bookTitle: values.sourceType === 2 ? values.bookTitle : undefined,
              bookPage: values.sourceType === 2 ? values.bookPage : undefined,
              instructions: values.sourceType === 3 ? values.instructions : undefined,
            }
           }
          },
          {
            onSuccess: (data) => {
              navigate(`/recipe/${data.createRecipe.id}`, {
                state: { successMessage: "Recipe created successfully!" }
              });
            }
          }
      );
    };

    return (
        <Page>
            <BackLink to="/recipe" />
            <PageTitle text={commonStrings.recipe.create} />
            {error && <Container><Alert message={extractErrorMessage(error)} closable /></Container>}
            <RecipeForm onSubmit={handleSubmit} isSubmitting={isPending} ingredients={ingredientsData?.ingredients || []} />
        </Page>
    );
}