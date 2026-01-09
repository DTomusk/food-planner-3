import { useQuery } from "@tanstack/react-query";
import { recipesQuery } from "../features/recipes/queries";
import PageTitle from "../components/PageTitle";
import RecipeList from "../features/recipes/components/RecipeList";
import Page from "../layout/PageWrapper";
import Spinner from "../components/Spinner";
import ErrorAlert from "../components/ErrorAlert";

export default function HomePage() {
// Use query subscribes to a cached query and updates the component when the data changes
  const {data, isLoading, error} = useQuery(recipesQuery);

  return (
    <Page>
    <PageTitle text="Home Page" />
    <h2 className="text-2xl font-bold underline">
      Recipes
    </h2>
    {isLoading && <Spinner/>}
    {error && <ErrorAlert message={(error as Error).message} />}
    {data && (
      <RecipeList recipes={data.recipes} />
    )}
  </Page>
  )
}