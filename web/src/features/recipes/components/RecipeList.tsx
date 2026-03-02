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
        <ul className="space-y-2">
        {recipes.map((recipe: {id: string, name: string}) => (
          <li key={recipe.id}>
            <RecipeListingCard recipe={recipe} 
              onClick={() => onCardClick ? onCardClick(recipe) : navigate(`/recipes/${recipe.id}`)}
              actions={renderActions ? renderActions(recipe) : undefined}/>
          </li>
        ))}
      </ul>
    )
}