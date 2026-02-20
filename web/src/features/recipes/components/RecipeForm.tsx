import { useFieldArray, useForm } from "react-hook-form";
import type { RecipeFormValues } from "../types";
import { Form, FormTitle, FormInputField, Button } from "@/components/";
import { commonStrings } from "@/lib/strings";

type RecipeFormProps = {
  onSubmit: (values: RecipeFormValues) => void;
  isSubmitting?: boolean;
  ingredients: { id: string; name: string }[];
};

export default function RecipeForm({
  onSubmit,
  isSubmitting = false,
  ingredients,
}: RecipeFormProps) {
  const {
    register,
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<RecipeFormValues>({
    defaultValues: {
      name: "",
      ingredientUsages: [{ ingredientId: "", quantity: 0 }],
    },
  });

  const { fields, append, remove } = useFieldArray({
    name: "ingredientUsages",
    control,
  });

  return (
    <Form onSubmit={handleSubmit(onSubmit)}>
      <FormTitle text="Create Recipe" />
      <FormInputField
        label="Recipe name"
        id="name"
        register={register("name", { required: "Name is required" })}
        error={errors.name}
        placeholder="Recipe name"
      />

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

      <Button disabled={isSubmitting} type="submit" loading={isSubmitting}>
        {commonStrings.forms.create}
      </Button>
    </Form>
  );
}
