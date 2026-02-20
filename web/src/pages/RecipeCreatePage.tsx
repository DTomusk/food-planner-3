import { Alert, BackLink, PageTitle } from "@/components";
import { useIngredients } from "@/features/ingredients/hooks/useIngredients";
import { RecipeForm, useCreateRecipe, type RecipeFormValues } from "@/features/recipes";
import { Page } from "@/layout";

export default function RecipeCreatePage() {
    const { mutate, isPending, error } = useCreateRecipe();
    const { data: ingredientsData } = useIngredients();

    const handleSubmit = (values: RecipeFormValues) => {
      mutate(
        // TODO: add ingredient usages
          { input: { 
            name: values.name, 
            ingredientUsages: values.ingredientUsages.map(usage => ({
              ingredientID: usage.ingredientId,
              quantity: usage.quantity,
              // Hardcode unit for now
              unit: 1
            }))
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
            <PageTitle text="Create Recipe" />
            {error && <Alert message={(error as Error).message} />}
            <RecipeForm onSubmit={handleSubmit} isSubmitting={isPending} ingredients={ingredientsData?.ingredients || []} />
        </Page>
    );
}