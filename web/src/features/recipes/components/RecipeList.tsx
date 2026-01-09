import type { Recipe } from "../types"
import RecipeCard from "./RecipeCard"

interface RecipeListProps {
    recipes: Array<Recipe>
}

export default function RecipeList({ recipes }: RecipeListProps) {
    return (
        <ul className="space-y-2">
        {recipes.map((recipe: {id: string, name: string}) => (
          <li key={recipe.id}>
            <RecipeCard recipe={recipe} />
          </li>
        ))}
      </ul>
    )
}