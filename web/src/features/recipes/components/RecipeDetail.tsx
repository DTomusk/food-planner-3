import Column from "@/components/ui/Column";
import { recipeStrings } from "../strings";
import type { Recipe } from "../types"

interface RecipeDetailProps {
    recipe: Recipe
}

export default function RecipeDetail({ recipe }: RecipeDetailProps) {
    if (!recipe) {
        return <p>{recipeStrings.errors.noRecipesFound}</p>;
    }

    return (
      <Column>
      <p className="text-center">Shared by {recipe.user}</p>
      <ul>
        {recipe.ingredients.map((ingredient, index) => (
          <li key={index}>{ingredient}</li>
        ))}
      </ul>
      </Column>
    )
}