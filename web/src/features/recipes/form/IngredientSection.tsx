import { useFieldArray, useFormContext } from "react-hook-form";
import type { RecipeFormValues } from "../types";
import Button from "@/components/ui/Button";
import FormSection from "@/components/form/FormSection";
import IngredientSelector from "./IngredientSelector";
import { useAvailableIngredients } from "../hooks/useAvailableIngredients";

interface IngredientSectionProps {
  ingredients: { id: string; name: string, counter: string; preferredUnit: { val: number; name: string; symbol: string } }[];
};

export default function IngredientSection({ ingredients }: IngredientSectionProps) {
    const { control, register, formState: { errors }, setValue } = useFormContext<RecipeFormValues>();

    const { fields, append, remove } = useFieldArray({
        name: "ingredientUsages",
        control,
    });

    const { getAvailableIngredients } = useAvailableIngredients(control, ingredients);
    
    return (
        <FormSection title="Ingredients" collapsible defaultCollapsed={false}>
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