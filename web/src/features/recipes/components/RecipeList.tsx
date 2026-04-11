import { recipeStrings } from "../strings";
import type { RecipeSummary } from "../types"
import { useNavigate } from "react-router-dom"
import RecipeListingCard from "./RecipeListingCard";
import { Alert } from "@/components";

interface RecipeListProps {
    recipes: Array<RecipeSummary>;
    renderActions?: (recipe: RecipeSummary) => React.ReactNode;
    onCardClick?: (recipe: RecipeSummary) => void;
}

export default function RecipeList({ recipes, renderActions, onCardClick }: RecipeListProps) {
    const navigate = useNavigate();

    if (recipes.length === 0) {
        return <Alert message={recipeStrings.errors.noRecipesFound} type="info" />;
    }

    return (
      <ul className="w-full space-y-8">
        {recipes.map((recipe: RecipeSummary) => (
          <li key={recipe.id} className="mx-auto w-full max-w-md md:max-w-none">
            <RecipeListingCard recipe={recipe} 
              onClick={() => onCardClick ? onCardClick(recipe) : navigate(`/recipes/${recipe.id}`)}
              actions={renderActions ? renderActions(recipe) : undefined}
              dietLevel={recipe.animalProductLevel}
            />
          </li>
        ))}
      </ul>
    )
}