import { Alert, BackLink, PageTitle } from "@/components";
import { useIngredients } from "@/features/ingredients/hooks/useIngredients";
import { RecipeForm, useCreateRecipe, type RecipeFormValues } from "@/features/recipes";
import { Page } from "@/layout";
import { commonStrings } from "@/lib/strings";

export default function RecipeCreatePage() {
    const { mutate, isPending, error } = useCreateRecipe();
    const { data: ingredientsData } = useIngredients();

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
            portions: values.portions
           }
          },
          {
            onSuccess: () => {
              alert("Recipe created successfully!");
            }
          }
      );
    };

    return (
        <Page>
            <BackLink to="/recipe" />
            <PageTitle text={commonStrings.recipe.create} />
            {error && <Alert message={(error as Error).message} />}
            <RecipeForm onSubmit={handleSubmit} isSubmitting={isPending} ingredients={ingredientsData?.ingredients || []} />
        </Page>
    );
}