import { useQuery } from "@tanstack/react-query";
import { graphqlClient } from "./lib/graphqlClient"

function App() {
  const getRecipes = `
    query GetRecipes {
    recipes {
      name
      id
    }
  }`;

  // Use query subscribes to a cached query and updates the component when the data changes
  // The query key is used to identify the query in the cache
  // The query function produces a promise that resolves to the data
  const {data, isLoading, error} = useQuery({
    queryKey: ['recipes'],
    queryFn: async () => graphqlClient.request(getRecipes),
  });
  console.log(data);

  return (
    <>
    <h1 className="text-3xl font-bold underline">
      Hello World!
    </h1>
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

export default App
