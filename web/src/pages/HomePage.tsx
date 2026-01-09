import { useQuery } from "@tanstack/react-query";
import { recipesQuery } from "../features/recipes/queries";
import PageTitle from "../components/PageTitle";
import RecipeList from "../features/recipes/components/RecipeList";
import Page from "../layout/PageWrapper";
import Spinner from "../components/Spinner";
import Alert from "../components/Alert";
import Button from "../components/Button";
import { useNavigate } from "react-router-dom";

export default function HomePage() {
// Use query subscribes to a cached query and updates the component when the data changes
  const {data, isLoading, error} = useQuery(recipesQuery);
  const navigate = useNavigate();

  return (
    <Page>
    <PageTitle text="Home Page" />
    <h2 className="text-2xl font-bold underline">
      Recipes
    </h2>
    <Button text="Create New Recipe" onClick={() => navigate("/recipe/create")} />
    {isLoading && <Spinner/>}
    {error && <Alert message={(error as Error).message} />}
    {data && (
      <RecipeList recipes={data.recipes} />
    )}
  </Page>
  )
}