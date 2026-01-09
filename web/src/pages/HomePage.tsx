import { useQuery } from "@tanstack/react-query";
import { recipesQuery } from "../features/recipes/queries";
import Title from "../components/Title";

export default function HomePage() {
// Use query subscribes to a cached query and updates the component when the data changes
  const {data, isLoading, error} = useQuery(recipesQuery);

  return (
    <>
    <Title text="Home Page" />
    <h2 className="text-2xl font-bold underline">
      Recipes
    </h2>
    {isLoading && <p>Loading...</p>}
    {error && <p>Error: {(error as Error).message}</p>}
    {data && (
      <ul>
        {data.recipes.map((recipe: {id: string, name: string}) => (
          <li key={recipe.id}>{recipe.name}</li>
        ))}
      </ul>
    )}
  </>
  )
}