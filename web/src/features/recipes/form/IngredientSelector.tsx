import { useWatch, type Control, type UseFormRegister, type UseFormSetValue } from "react-hook-form";
import type { RecipeFormValues } from "../types";
import Input from "@/components/ui/Input";
import FormField from "@/components/form/FormField";
import Select from "@/components/ui/Select";
import IconButton from "@/components/ui/IconButton";
import { X } from "lucide-react";
import Inline from "@/components/layout/Inline";
import { useEffect } from "react";

interface IngredientSelectorProps {
  index: number;
  control: Control<RecipeFormValues>;
  register: UseFormRegister<RecipeFormValues>;
  ingredients: { id: string; name: string; counter: string; preferredUnit: { val: number; name: string; symbol: string } }[];
  remove: (index: number) => void;
  canRemove: boolean;
  errors: any;
  setValue: UseFormSetValue<RecipeFormValues>;
};

export default function IngredientSelector({ index, control, register, ingredients, remove, canRemove, errors, setValue }: IngredientSelectorProps) {
    const selectedIngredient = useWatch({
        control,
        name: `ingredientUsages.${index}.ingredientId`,
    });

    // TODO: move somewhere
    const formatIngredientOption = (ingredientName: string, counter: string) => {
      if (!counter) return ingredientName;
      return `${ingredientName} (${counter})`;
    }

    // Get which ingredient is selected to determine if we should show unit
    const selectedIngredientObj = ingredients.find((ingredient) => ingredient.id === selectedIngredient);

    // If quantum, then unitless, so no unit select
    // This is a bit brittle, we could send a flag with the unit to show whether it should be displayed 
    // But the unitless unit might be the only time we have to do this, so let's keep it like this and see if that changes
    const isQuantumUnit = selectedIngredientObj?.preferredUnit.val === 1;

    const hasSelection = Boolean(selectedIngredient);

    useEffect(() => {
      if (selectedIngredientObj) {
        setValue(
          `ingredientUsages.${index}.unit`,
          selectedIngredientObj.preferredUnit.val
        );
      }
    }, [selectedIngredientObj, index, setValue]);

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
                {formatIngredientOption(ingredient.name, ingredient.counter)}
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

      {selectedIngredientObj && !isQuantumUnit && (
        <FormField htmlFor={`ingredientUsages.${index}.unit`} label="Unit">
          <Select disabled value={selectedIngredientObj.preferredUnit.val}>
            <option value={selectedIngredientObj.preferredUnit.val}>{selectedIngredientObj.preferredUnit.name}</option>
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