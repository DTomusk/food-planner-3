import { useWatch, type Control, type UseFormRegister } from "react-hook-form";
import type { RecipeFormValues } from "../types";
import Input from "@/components/ui/Input";
import FormField from "@/components/form/FormField";
import Select from "@/components/ui/Select";
import IconButton from "@/components/ui/IconButton";
import { X } from "lucide-react";
import Inline from "@/components/layout/Inline";

interface IngredientSelectorProps {
  index: number;
  control: Control<RecipeFormValues>;
  register: UseFormRegister<RecipeFormValues>;
  ingredients: { id: string; name: string; preferredUnit: { val: number; name: string; symbol: string } }[];
  remove: (index: number) => void;
  canRemove: boolean;
  errors: any;
};

export default function IngredientSelector({ index, control, register, ingredients, remove, canRemove, errors }: IngredientSelectorProps) {
    const selectedIngredient = useWatch({
        control,
        name: `ingredientUsages.${index}.ingredientId`,
    });

    // Get which ingredient is selected to determine if we should show unit
    const selectedIngredientObj = ingredients.find((ingredient) => ingredient.id === selectedIngredient);

    // If quantum, then unitless, so no unit select
    // This is a bit brittle, we could send a flag with the unit to show whether it should be displayed 
    // But the unitless unit might be the only time we have to do this, so let's keep it like this and see if that changes
    const isQuantumUnit = selectedIngredientObj?.preferredUnit.val === 1;

    const hasSelection = Boolean(selectedIngredient);

    return (
    <Inline>
      <div className="flex-1">
        <FormField htmlFor={`ingredientUsages.${index}.ingredientId`} label="Ingredient" error={errors.ingredientUsages?.[index]?.ingredientId?.message}>
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
              validate: (value) => value > 0 || "Must be positive",
            })}
          />
        </FormField>
      )}

      {selectedIngredientObj && !isQuantumUnit && (
        <FormField htmlFor={`ingredientUsages.${index}.unit`} label="Unit">
          <Select disabled value={selectedIngredientObj.preferredUnit.val}>
            <option value={selectedIngredientObj.preferredUnit.val}>{selectedIngredientObj.preferredUnit.name}</option>
          </Select>
        </FormField>
      )}

      {canRemove && (
        <IconButton type="button" variant="secondary" onClick={() => remove(index)} aria-label="Remove ingredient">
          <X size={20} />
        </IconButton>
      )}
    </Inline>
    );
}