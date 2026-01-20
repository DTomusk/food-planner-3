import { Alert, BackLink, PageTitle } from "@/components";
import { RecipeForm, useCreateRecipe, type RecipeFormValues } from "@/features/recipes";
import { Page } from "@/layout";

export default function RecipeCreatePage() {
    const { mutate, isPending, error } = useCreateRecipe();

    const handleSubmit = (values: RecipeFormValues) => {
      mutate(
          { input: { name: values.name } },
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
            <RecipeForm onSubmit={handleSubmit} isSubmitting={isPending} />
        </Page>
    );
}