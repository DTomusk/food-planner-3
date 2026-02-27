import { useFieldArray, type Control, type FieldErrors, type UseFormRegister, type UseFormSetValue } from "react-hook-form";
import type { RecipeFormValues } from "../types";
import Button from "@/components/ui/Button";
import FormSection from "@/components/form/FormSection";
import IngredientSelector from "./IngredientSelector";
import { useAvailableIngredients } from "../hooks/useAvailableIngredients";

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

    const { getAvailableIngredients } = useAvailableIngredients(control, ingredients);
    
    return (
        <FormSection title="Ingredients">
            {fields.map((field, index) => (
                <IngredientSelector 
                    key={field.id}
                    index={index}
                    control={control}
                    register={register}
                    ingredients={getAvailableIngredients(index)}
                    remove={remove}
                    canRemove={fields.length > 1}
                    errors={errors}
                    setValue={setValue}
                />
            ))}

            <Button type="button" onClick={() => append({ ingredientId: "", quantity: 0, unit: 1 })}>
                Add Ingredient
            </Button>
        </FormSection>
    );
}