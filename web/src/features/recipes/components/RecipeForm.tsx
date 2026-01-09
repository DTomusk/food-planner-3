import { useForm } from "react-hook-form";
import type { RecipeFormValues } from "../types";

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
    <form onSubmit={handleSubmit(onSubmit)}>
      <div>
        <input
          {...register("name", { required: "Name is required" })}
          placeholder="Recipe name"
        />
        {errors.name && <p>{errors.name.message}</p>}
      </div>

      <button type="submit" disabled={isSubmitting}>
        Create
      </button>
    </form>
  );
}
