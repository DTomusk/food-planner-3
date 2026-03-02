import { useMemo } from "react";
import { useWatch, type Control } from "react-hook-form";
import type { RecipeFormValues } from "../types";

interface Ingredient {
  id: string;
  name: string;
  counter: string;
  preferredUnit: {
    val: number;
    name: string;
    symbol: string;
  };
}

export function useAvailableIngredients(
  control: Control<RecipeFormValues>,
  allIngredients: Ingredient[]
) {
  // Watch which ingredients have been selected to later filter for duplicates
  const ingredientUsages = useWatch({
    control,
    name: "ingredientUsages",
  });

  // Cache ids to prevent recalculating on every render
  const selectedIngredientIds = useMemo(() => {
    return (
      ingredientUsages
        ?.map((usage) => usage?.ingredientId)
        .filter(Boolean) ?? []
    );
  }, [ingredientUsages]);

  // Function to get allowed ingredients for a given row
  function getAvailableIngredients(index: number) {
    const currentSelection = ingredientUsages?.[index]?.ingredientId;

    return allIngredients.filter((ingredient) => {
      if (ingredient.id === currentSelection) return true;
      return !selectedIngredientIds.includes(ingredient.id);
    });
  }

  return {
    ingredientUsages,
    getAvailableIngredients,
  };
}