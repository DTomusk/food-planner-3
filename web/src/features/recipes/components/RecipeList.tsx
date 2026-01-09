import type { Recipe } from "../types"

interface RecipeListProps {
    recipes: Array<Recipe>
}

export default function RecipeList({ recipes }: RecipeListProps) {
    return (
        <ul>
        {recipes.map((recipe: {id: string, name: string}) => (
          <li key={recipe.id}>{recipe.name}</li>
        ))}
      </ul>
    )
}