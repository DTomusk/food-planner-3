import { recipeStrings } from "../strings";
import type { Recipe } from "../types"
import RecipeCard from "./RecipeCard"
import { useNavigate } from "react-router-dom"

interface RecipeListProps {
    recipes: Array<Recipe>
}

export default function RecipeList({ recipes }: RecipeListProps) {
    const navigate = useNavigate();

    if (recipes.length === 0) {
        return <p>{recipeStrings.errors.noRecipesFound}</p>;
    }

    return (
        <ul className="space-y-2">
        {recipes.map((recipe: {id: string, name: string}) => (
          <li key={recipe.id}>
            <RecipeCard recipe={recipe} onClick={() => navigate(`/recipe/${recipe.id}`)}/>
          </li>
        ))}
      </ul>
    )
}