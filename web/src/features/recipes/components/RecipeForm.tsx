import { useForm } from "react-hook-form";
import type { RecipeFormValues } from "../types";
import { Form, FormTitle, FormInputField, Button } from "@/components/";
import { commonStrings } from "@/lib/strings";
import IngredientSection from "./IngredientSection";
import FormSection from "@/components/FormSection";

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

  return (
    <Form onSubmit={handleSubmit(onSubmit)}>
      <FormTitle text="Create Recipe" />
      <FormSection title="Recipe Details">
        <FormInputField
          label="Recipe name"
          id="name"
          register={register("name", { required: "Name is required" })}
          error={errors.name}
          placeholder="Recipe name"
        />
      </FormSection>      

      <IngredientSection control={control} register={register} errors={errors} ingredients={ingredients} />

      <Button disabled={isSubmitting} type="submit" loading={isSubmitting}>
        {commonStrings.forms.create}
      </Button>
    </Form>
  );
}
