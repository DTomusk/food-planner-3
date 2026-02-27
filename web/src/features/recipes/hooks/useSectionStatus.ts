import { useFormContext } from "react-hook-form";
import type { RecipeFormValues } from "../types";

export function useSectionStatus(fields: readonly (keyof RecipeFormValues)[]) {
  const { getValues, formState: { errors } } = useFormContext<RecipeFormValues>();
  
  const hasError = fields.some(field => errors[field]);
  if (hasError) return "error";

  const isComplete = fields.every(field => {
    const val = getValues(field);
    return val !== "" && val !== 0 && val != null;
  });

  return isComplete ? "completed" : "pending";
}