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
import SearchDropdown from "@/components/ui/SearchDropdown";
import Fuse from "fuse.js";
import { useMemo } from "react";

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

  const dropdownItems = useMemo(() =>
    ingredients.map((ingredient) => ({
    label: formatIngredientOption(ingredient),
    value: ingredient.id,
  })), [ingredients, formatIngredientOption]);

  const fuse = useMemo(() => new Fuse(
    dropdownItems,
    {
      keys: ["label"],
      threshold: 0.3,
    }
  ), [dropdownItems]);
  
    return (
    <Inline align="center">
      <div className="flex-1">
        <FormField htmlFor={`ingredientUsages.${index}.ingredientId`} label="Ingredient" error={errors.ingredientUsages?.[index]?.ingredientId?.message}>
          <SearchDropdown
            maxResults={5}
            items={dropdownItems}
            onSelect={(item) => {
              setValue(`ingredientUsages.${index}.ingredientId`, item?.value ?? "", { shouldDirty: true });
            }}
            selectedItem={selectedIngredient ? { label: formatIngredientOption(selectedIngredient), value: selectedIngredient.id } : null}
            filterFunction={(query) => query ? fuse.search(query).map(result => result.item) : dropdownItems}
          />
        </FormField>
      </div>

      <FormField htmlFor={`ingredientUsages.${index}.quantity`} label="Quantity" error={errors.ingredientUsages?.[index]?.quantity?.message}>
        <Input
          type="number"
          min={0}
          step={1}
          placeholder="Qty"
          {...register(`ingredientUsages.${index}.quantity`, {
            valueAsNumber: true})}
          disabled={!hasSelection}
        />
      </FormField>

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