import { useFieldArray, type Control, type FieldErrors, type UseFormRegister } from "react-hook-form";
import type { RecipeFormValues } from "../types";
import Button from "@/components/ui/Button";
import FormSection from "@/components/form/FormSection";
import IngredientSelector from "./IngredientSelector";

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
                <IngredientSelector 
                    key={field.id}
                    index={index}
                    control={control}
                    register={register}
                    ingredients={ingredients}
                    remove={remove}
                    canRemove={fields.length > 1}
                />
            ))}

            <Button type="button" onClick={() => append({ ingredientId: "", quantity: 0 })}>
                Add Ingredient
            </Button>
        </FormSection>
    );
}