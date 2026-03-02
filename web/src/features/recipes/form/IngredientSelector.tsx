import { type Control, type UseFormRegister, type UseFormSetValue } from "react-hook-form";
import type { RecipeFormValues } from "../types";
import Input from "@/components/ui/Input";
import FormField from "@/components/form/FormField";
import Select from "@/components/ui/Select";
import IconButton from "@/components/ui/IconButton";
import { X } from "lucide-react";
import Inline from "@/components/layout/Inline";
import { useIngredientSelector } from "../hooks/useIngredientSelector";
import type { IngredientOptionModel } from "@/features/ingredients/types";

interface IngredientSelectorProps {
  index: number;
  control: Control<RecipeFormValues>;
  register: UseFormRegister<RecipeFormValues>;
  ingredients: IngredientOptionModel[];
  remove: (index: number) => void;
  canRemove: boolean;
  errors: any;
  setValue: UseFormSetValue<RecipeFormValues>;
};

export default function IngredientSelector({ index, control, register, ingredients, remove, canRemove, errors, setValue }: IngredientSelectorProps) {
    const {
    selectedIngredient,
    hasSelection,
    isQuantumUnit,
    formatIngredientOption,
  } = useIngredientSelector({
    index,
    control,
    setValue,
    ingredients,
  });  
  
    return (
    <Inline>
      <div className="flex-1">
        <FormField htmlFor={`ingredientUsages.${index}.ingredientId`} label="Ingredient" error={errors.ingredientUsages?.[index]?.ingredientId?.message}>
          <Select defaultValue="" disabled={ingredients.length === 0} {...register(`ingredientUsages.${index}.ingredientId`)}>
            <option value="" disabled={hasSelection}>
              Select ingredient
            </option>
            {ingredients.map((ingredient) => (
              <option key={ingredient.id} value={ingredient.id}>
                {formatIngredientOption(ingredient)}
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
              valueAsNumber: true})}
          />
        </FormField>
      )}

      {selectedIngredient && !isQuantumUnit && (
        <FormField htmlFor={`ingredientUsages.${index}.unit`} label="Unit">
          <Select disabled value={selectedIngredient.preferredUnit.val}>
            <option value={selectedIngredient.preferredUnit.val}>
            {selectedIngredient.preferredUnit.name}
          </option>          
          </Select>
        </FormField>
      )}

      <input
        type="hidden"
        {...register(`ingredientUsages.${index}.unit`, {
          valueAsNumber: true,
        })}
      />

      {canRemove && (
        <IconButton type="button" variant="secondary" onClick={() => remove(index)} aria-label="Remove ingredient">
          <X size={20} />
        </IconButton>
      )}
    </Inline>
    );
}