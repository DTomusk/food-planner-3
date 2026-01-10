import { useForm } from "react-hook-form";
import type { RecipeFormValues } from "../types";
import { Form, FormTitle, FormInputField, FormSubmitButton } from "@/components/";

type RecipeFormProps = {
  onSubmit: (values: RecipeFormValues) => void;
  isSubmitting?: boolean;
};

export default function RecipeForm({
  onSubmit,
  isSubmitting = false,
}: RecipeFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RecipeFormValues>();

  return (
    <Form onSubmit={handleSubmit(onSubmit)}>
      <FormTitle text="Create Recipe" />
      <FormInputField
        label="Recipe Name"
        id="name"
        register={register("name", { required: "Name is required" })}
        error={errors.name}
        placeholder="Recipe name"
      />
      <FormSubmitButton disabled={isSubmitting}>
        Create
      </FormSubmitButton>
    </Form>
  );
}
