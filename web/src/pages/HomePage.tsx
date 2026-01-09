import PageTitle from "../components/PageTitle";
import RecipeList from "../features/recipes/components/RecipeList";
import Page from "../layout/PageWrapper";
import Spinner from "../components/Spinner";
import Alert from "../components/Alert";
import { useRecipes } from "../features/recipes/hooks/useRecipes";
import { RecipeForm } from "../features/recipes/components/RecipeForm";
import type { RecipeFormValues } from "../features/recipes/types";
import { useCreateRecipe } from "../features/recipes/hooks/useCreateRecipe";

export default function HomePage() {
  const {data, isLoading, error: fetchError} = useRecipes();
  const { mutate, isPending, error: createError } = useCreateRecipe();

    const handleSubmit = (values: RecipeFormValues) => {
        mutate(
            { input: { name: values.name } },
            {
              onSuccess: () => {
                alert("Recipe created successfully!");
              }
            }
        );
    }
  return (
    <Page>
    <PageTitle text="Home Page" />
    <h2 className="text-2xl font-bold underline">
      Recipes
    </h2>
    <RecipeForm
            onSubmit={handleSubmit}
            isSubmitting={isPending}
        />
    {isLoading && <Spinner/>}
    {fetchError && <Alert message={(fetchError as Error).message} />}
    {createError && <Alert message={(createError as Error).message} />}
    {data && (
      <RecipeList recipes={data.recipes} />
    )}
  </Page>
  )
}