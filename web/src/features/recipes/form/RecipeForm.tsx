import { FormProvider, useForm } from "react-hook-form";
import type { RecipeFormValues } from "../types";
import { Form, Button } from "@/components/";
import { commonStrings } from "@/lib/strings";
import IngredientSection from "./IngredientSection";
import Stack from "@/components/layout/Stack";
import RecipeSourceSection from "./RecipeSourceSection";
import RecipeDetailSection from "./RecipeDetailSection";
import Container from "@/components/layout/Container";

type RecipeFormProps = {
  onSubmit: (values: RecipeFormValues) => void;
  isSubmitting?: boolean;
  ingredients: { id: string; name: string, preferredUnit: { val: number; name: string; symbol: string } }[];
};

export default function RecipeForm({
  onSubmit,
  isSubmitting = false,
  ingredients,
}: RecipeFormProps) {

  const methods = useForm<RecipeFormValues>({
    mode: "onChange",
    defaultValues: {
      name: "",
      ingredientUsages: [{ ingredientId: "", quantity: 0, unit: 1 }],
      prepMins: 0,
      cookMins: 0,
      portions: 0,
      sourceType: 0
    },
  });

  const {
    handleSubmit, 
    formState: { isValid, isSubmitted },
  } = methods;

  return (
    <Container size="md">
      <FormProvider {...methods}>
      <Form onSubmit={handleSubmit(onSubmit)}>
        <Stack space="lg">
          <RecipeDetailSection />
          <IngredientSection ingredients={ingredients} />
          <RecipeSourceSection />
          <Button disabled={isSubmitting || (isSubmitted && !isValid)} type="submit" loading={isSubmitting}>
            {commonStrings.forms.create}
          </Button>
        </Stack>
      </Form>
      </FormProvider>
    </Container>
  );
}
