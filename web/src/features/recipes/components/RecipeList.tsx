import { recipeStrings } from "../strings";
import type { RecipeSummary } from "../types"
import { useNavigate } from "react-router-dom"
import RecipeListingCard from "./RecipeListingCard";
import { Alert } from "@/components";

interface RecipeListProps {
    recipes: Array<RecipeSummary>;
}

export default function RecipeList({ recipes }: RecipeListProps) {
    const navigate = useNavigate();

    if (recipes.length === 0) {
        return <Alert message={recipeStrings.errors.noRecipesFound} type="info" />;
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