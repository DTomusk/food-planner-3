import { useEffect, useMemo } from "react";
import { useWatch, type Control, type UseFormSetValue } from "react-hook-form";
import type { RecipeFormValues } from "../types";
import type { IngredientOptionModel } from "@/features/ingredients/types";

interface UseIngredientSelectorProps {
  index: number;
  control: Control<RecipeFormValues>;
  setValue: UseFormSetValue<RecipeFormValues>;
  ingredients: IngredientOptionModel[];
}

export function useIngredientSelector({
  index,
  control,
  setValue,
  ingredients,
}: UseIngredientSelectorProps) {
  const selectedIngredientId = useWatch({
    control,
    name: `ingredientUsages.${index}.ingredientId`,
  });

  const selectedIngredient = useMemo(
    () => ingredients.find((i) => i.id === selectedIngredientId),
    [ingredients, selectedIngredientId]
  );

  const hasSelection = Boolean(selectedIngredientId);

  const isQuantumUnit =
    selectedIngredient?.preferredUnit.val === 1;

  // keep unit in sync when ingredient changes
  useEffect(() => {
    if (selectedIngredient) {
      setValue(
        `ingredientUsages.${index}.unit`,
        selectedIngredient.preferredUnit.val
      );
    }
  }, [selectedIngredient, index, setValue]);

  const formatIngredientOption = (ingredient: IngredientOptionModel) => {
    if (!ingredient.counter) return ingredient.name;
    return `${ingredient.name} (${ingredient.counter})`;
  };

  return {
    selectedIngredient,
    hasSelection,
    isQuantumUnit,
    formatIngredientOption,
  };
}