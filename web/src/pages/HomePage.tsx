import { Alert, PageTitle, Spinner } from "@/components";
import { RecipeList, useRecipes, useCreateRecipe, type RecipeFormValues, RecipeForm} from "@/features/recipes";
import { Page } from "@/layout";

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
    <RecipeForm
            onSubmit={handleSubmit}
            isSubmitting={isPending}
        />
    <h2 className="text-2xl font-bold underline">
      Recipes
    </h2>
    {isLoading && <Spinner/>}
    {fetchError && <Alert message={(fetchError as Error).message} />}
    {createError && <Alert message={(createError as Error).message} />}
    {data && (
      <RecipeList recipes={data.recipes} />
    )}
  </Page>
  )
}