import { FormProvider, useForm } from "react-hook-form";
import type { RecipeFormValues } from "../types";
import { Form, Button, Inline } from "@/components/";
import { commonStrings } from "@/lib/strings";
import IngredientSection from "./IngredientSection";
import Stack from "@/components/layout/Stack";
import RecipeSourceSection from "./RecipeSourceSection";
import RecipeDetailSection from "./RecipeDetailSection";
import Container from "@/components/layout/Container";
import { recipeFormSchema } from "../schemas/recipeFormSchema";
import { zodResolver } from "@hookform/resolvers/zod";
import type { IngredientOptionModel } from "@/features/ingredients/types";
import { DEFAULT_RECIPE_FORM_VALUES } from "../mappers/recipeFormMapper";

type RecipeFormProps = {
  onSubmit: (values: RecipeFormValues) => void;
  isSubmitting?: boolean;
  ingredients: IngredientOptionModel[];
  defaultValues?: RecipeFormValues;
};

export default function RecipeForm({
  onSubmit,
  isSubmitting = false,
  ingredients,
  defaultValues,
}: RecipeFormProps) {
  const methods = useForm<RecipeFormValues>({
    resolver: zodResolver(recipeFormSchema),
    mode: "onChange",
    defaultValues: defaultValues ?? DEFAULT_RECIPE_FORM_VALUES,
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
            <Inline justify="start">
            <Button disabled={isSubmitting || (isSubmitted && !isValid)} type="submit" loading={isSubmitting}>
              {commonStrings.forms.create}
            </Button>
            </Inline>
          </Stack>
        </Form>
      </FormProvider>
    </Container>
  );
}
