import { useForm } from "react-hook-form";
import type { RecipeFormValues } from "../types";
import Form from "../../../components/FormWrapper";
import { InputField } from "../../../components/FormInputField";
import SubmitButton from "../../../components/FormSubmitButton";

type RecipeFormProps = {
  onSubmit: (values: RecipeFormValues) => void;
  isSubmitting?: boolean;
};

export function RecipeForm({
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
      <InputField
        label="Recipe Name"
        id="name"
        register={register("name", { required: "Name is required" })}
        error={errors.name}
        placeholder="Recipe name"
      />
      <SubmitButton disabled={isSubmitting}>
        Create
      </SubmitButton>
    </Form>
  );
}
