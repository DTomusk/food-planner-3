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
      <></>
    )
}