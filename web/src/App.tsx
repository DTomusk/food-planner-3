import { queryClient } from "./app/queryClient"
import { graphqlClient } from "./lib/graphqlClient"

function App() {
  const getRecipes = `
  query GetRecipes {
  recipes {
    name
    id
  }
}`

async function fetchRecipes() {
  const response = await queryClient.fetchQuery({
    queryKey: ['recipes'],
    queryFn: async () => graphqlClient.request(getRecipes),
  });
  console.log(response);
}

fetchRecipes();

  return (
    <h1 className="text-3xl font-bold underline">
      Hello World!
    </h1>
  )
}

export default App
