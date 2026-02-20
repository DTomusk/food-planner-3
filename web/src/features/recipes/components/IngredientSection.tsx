import { useFieldArray, type Control, type FieldErrors, type UseFormRegister } from "react-hook-form";
import type { RecipeFormValues } from "../types";
import Button from "@/components/Button";
import FormSection from "@/components/FormSection";

interface IngredientSectionProps {
  control: Control<RecipeFormValues>;
  register: UseFormRegister<RecipeFormValues>;
  errors: FieldErrors<RecipeFormValues>;
  ingredients: { id: string; name: string }[];
};

export default function IngredientSection({ control, register, errors, ingredients }: IngredientSectionProps) {
    const { fields, append, remove } = useFieldArray({
        name: "ingredientUsages",
        control,
    });
    
    return (
        <FormSection title="Ingredients">
            {fields.map((field, index) => (
                <div key={field.id}>
                <select {...register(`ingredientUsages.${index}.ingredientId`, { required: "Ingredient is required" })}>
                    <option value="">Select ingredient</option>
                    {ingredients.map((ingredient) => (
                    <option key={ingredient.id} value={ingredient.id}>
                        {ingredient.name}
                    </option>
                    ))}
                </select>
                <input
                    type="number"
                    {...register(`ingredientUsages.${index}.quantity`, { required: "Quantity is required", min: 0 })}
                />

                {fields.length > 1 && (
                    <Button type="button" onClick={() => remove(index)}>
                    Remove
                    </Button>
                )}
                </div>
            ))}

            <Button type="button" onClick={() => append({ ingredientId: "", quantity: 0 })}>
                Add Ingredient
            </Button>
        </FormSection>
    );
}