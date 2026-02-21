import { useWatch, type Control, type UseFormRegister } from "react-hook-form";
import type { RecipeFormValues } from "../types";
import Input from "@/components/ui/Input";

interface IngredientSelectorProps {
  index: number;
  control: Control<RecipeFormValues>;
  register: UseFormRegister<RecipeFormValues>;
  ingredients: { id: string; name: string }[];
  remove: (index: number) => void;
  canRemove: boolean;
};

export default function IngredientSelector({ index, control, register, ingredients, remove, canRemove }: IngredientSelectorProps) {
    const selectedIngredient = useWatch({
        control,
        name: `ingredientUsages.${index}.ingredientId`,
    });

    const hasSelection = Boolean(selectedIngredient);

    return (
    <div className="flex items-end gap-3 outline outline-1">
      <div className="flex-1">
        <select
          {...register(`ingredientUsages.${index}.ingredientId`, {
            required: "Ingredient is required",
          })}
        >
          <option value="" disabled={hasSelection}>
            Select ingredient
          </option>
          {ingredients.map((ingredient) => (
            <option key={ingredient.id} value={ingredient.id}>
              {ingredient.name}
            </option>
          ))}
        </select>
      </div>

      {hasSelection && (
        <div className="w-24">
          <Input
            type="number"
            min={0}
            step={1}
            placeholder="Qty"
            {...register(`ingredientUsages.${index}.quantity`, {
              required: "Quantity is required",
              valueAsNumber: true,
              min: { value: 0, message: "Must be positive" },
            })}
          />
        </div>
      )}

      {canRemove && (
        <button
          type="button"
          onClick={() => remove(index)}
          className="text-sm text-red-500"
        >
          Remove
        </button>
      )}
    </div>
    );
}