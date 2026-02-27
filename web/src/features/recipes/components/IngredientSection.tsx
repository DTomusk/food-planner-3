import { useFieldArray, useWatch, type Control, type FieldErrors, type UseFormRegister, type UseFormSetValue } from "react-hook-form";
import type { RecipeFormValues } from "../types";
import Button from "@/components/ui/Button";
import FormSection from "@/components/form/FormSection";
import IngredientSelector from "./IngredientSelector";

interface IngredientSectionProps {
  control: Control<RecipeFormValues>;
  register: UseFormRegister<RecipeFormValues>;
  errors: FieldErrors<RecipeFormValues>;
  setValue: UseFormSetValue<RecipeFormValues>;
  ingredients: { id: string; name: string, preferredUnit: { val: number; name: string; symbol: string } }[];
};

export default function IngredientSection({ control, register, errors, setValue, ingredients }: IngredientSectionProps) {
    const { fields, append, remove } = useFieldArray({
        name: "ingredientUsages",
        control,
    });

    // Watch which ingredients have been selected to later filter for duplicates
    const ingredientUsages = useWatch({
        control,
        name: "ingredientUsages",
    });

    const selectedIngredientIds = ingredientUsages?.map((usage) => usage?.ingredientId).filter(Boolean) || [];
    
    return (
        <FormSection title="Ingredients">
            {fields.map((field, index) => {
                // Get the ingredient that's currently selected in this select (that should still show up in the dropdown)
                const currentSelection = ingredientUsages?.[index]?.ingredientId;

                // Filter out ingredients that have been selected in other selectors
                const availableIngredients = ingredients.filter((ingredient) => {
                    if (ingredient.id === currentSelection) return true;
                    return !selectedIngredientIds.includes(ingredient.id);
                });

                return(
                <IngredientSelector 
                    key={field.id}
                    index={index}
                    control={control}
                    register={register}
                    ingredients={availableIngredients}
                    remove={remove}
                    canRemove={fields.length > 1}
                    errors={errors}
                    setValue={setValue}
                />
            )})}

            <Button type="button" onClick={() => append({ ingredientId: "", quantity: 0, unit: 1 })}>
                Add Ingredient
            </Button>
        </FormSection>
    );
}