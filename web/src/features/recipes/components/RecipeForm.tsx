import { useForm } from "react-hook-form";
import type { RecipeFormValues } from "../types";
import { Form, Button } from "@/components/";
import { commonStrings } from "@/lib/strings";
import IngredientSection from "./IngredientSection";
import FormSection from "@/components/form/FormSection";
import Input from "@/components/ui/Input";
import FormField from "@/components/form/FormField";
import Stack from "@/components/layout/Stack";

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
    formState: { errors, isValid },
  } = useForm<RecipeFormValues>({
    defaultValues: {
      name: "",
      ingredientUsages: [{ ingredientId: "", quantity: 0 }],
    },
  });

  return (
    <Form onSubmit={handleSubmit(onSubmit)}>
      <Stack space="md">
      <FormSection>
        <FormField htmlFor="name" label="Recipe name" error={errors.name?.message}>
          <Input type="text" placeholder="Recipe name" 
          {...register("name", 
          { required: "Please add a name",
              minLength: { value: 3, message: "Name must be at least 3 characters" },
           })} />
        </FormField>
      </FormSection>      

      <IngredientSection control={control} register={register} errors={errors} ingredients={ingredients} />

      <Button disabled={isSubmitting || !isValid} type="submit" loading={isSubmitting}>
        {commonStrings.forms.create}
      </Button>
      </Stack>
    </Form>
  );
}
