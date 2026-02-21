import { useWatch, type Control, type UseFormRegister } from "react-hook-form";
import type { RecipeFormValues } from "../types";
import Input from "@/components/ui/Input";
import FormField from "@/components/form/FormField";
import Select from "@/components/ui/Select";

interface IngredientSelectorProps {
  index: number;
  control: Control<RecipeFormValues>;
  register: UseFormRegister<RecipeFormValues>;
  ingredients: { id: string; name: string }[];
  remove: (index: number) => void;
  canRemove: boolean;
  errors: any;
};

export default function IngredientSelector({ index, control, register, ingredients, remove, canRemove, errors }: IngredientSelectorProps) {
    const selectedIngredient = useWatch({
        control,
        name: `ingredientUsages.${index}.ingredientId`,
    });

    const hasSelection = Boolean(selectedIngredient);

    return (
    <div className="flex items-end gap-3">
      <div className="flex-1">
        <FormField htmlFor={`ingredientUsages.${index}.ingredientId`} label="Choose ingredient" error={errors.ingredientUsages?.[index]?.ingredientId?.message}>
          <Select defaultValue="" disabled={ingredients.length === 0} {...register(`ingredientUsages.${index}.ingredientId`, { required: "Ingredient is required" })}>
            <option value="" disabled={hasSelection}>
              Select ingredient
            </option>
            {ingredients.map((ingredient) => (
              <option key={ingredient.id} value={ingredient.id}>
                {ingredient.name}
              </option>
            ))}
          </Select>
        </FormField>
      </div>

      {hasSelection && (
        <FormField htmlFor={`ingredientUsages.${index}.quantity`} label="Quantity" error={errors.ingredientUsages?.[index]?.quantity?.message}>
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
        </FormField>
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