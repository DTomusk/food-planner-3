import { recipeStrings } from "../strings";
import type { RecipeSummary } from "../types"
import { useNavigate } from "react-router-dom"
import RecipeListingCard from "./RecipeListingCard";

interface RecipeListProps {
    recipes: Array<RecipeSummary>;
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
            <RecipeListingCard recipe={recipe} onClick={() => navigate(`/recipe/${recipe.id}`)}/>
          </li>
        ))}
      </ul>
    )
}