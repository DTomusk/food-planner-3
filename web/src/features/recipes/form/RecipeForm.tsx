import { FormProvider, useForm } from "react-hook-form";
import { RecipeSourceType, type RecipeFormValues } from "../types";
import { Form, Button } from "@/components/";
import { commonStrings } from "@/lib/strings";
import IngredientSection from "./IngredientSection";
import Stack from "@/components/layout/Stack";
import RecipeSourceSection from "./RecipeSourceSection";
import RecipeDetailSection from "./RecipeDetailSection";
import Container from "@/components/layout/Container";
import { recipeFormSchema } from "../schemas/recipeFormSchema";
import { zodResolver } from "@hookform/resolvers/zod";

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
    resolver: zodResolver(recipeFormSchema),
    mode: "onChange",
    defaultValues: {
      name: "",
      prepMins: 0,
      cookMins: 0,
      portions: 0,
      ingredientUsages: [{ ingredientId: "", quantity: 0, unit: 1 }],
      sourceType: RecipeSourceType.None,
      url: "",
      bookTitle: "",
      bookPage: 0,
      instructions: "",
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
