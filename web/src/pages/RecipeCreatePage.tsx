import { Alert, BackLink, PageTitle } from "@/components";
import Container from "@/components/layout/Container";
import { useIngredients } from "@/features/ingredients/hooks/useIngredients";
import { RecipeForm, useCreateRecipe, type RecipeFormValues } from "@/features/recipes";
import { RecipeSourceType } from "@/features/recipes/types";
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
              type: values.sourceType === RecipeSourceType.None ? 0 : values.sourceType === RecipeSourceType.Website ? 1 : values.sourceType === RecipeSourceType.Cookbook ? 2 : 3,
              url: values.sourceType === RecipeSourceType.Website ? values.url : undefined,
              bookTitle: values.sourceType === RecipeSourceType.Cookbook ? values.bookTitle : undefined,
              bookPage: values.sourceType === RecipeSourceType.Cookbook ? values.bookPage : undefined,
              instructions: values.sourceType === RecipeSourceType.Original ? values.instructions : undefined,
            }
           }
          },
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
            <RecipeForm onSubmit={handleSubmit} isSubmitting={isPending} ingredients={ingredientsData?.ingredients || []} />
        </Page>
    );
}