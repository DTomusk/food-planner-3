import { useEffect } from "react";
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
  isPreparingImageUpload?: boolean;
  ingredients: IngredientOptionModel[];
  onDirtyChange?: (isDirty: boolean) => void;
  defaultValues?: RecipeFormValues;
  imageFile?: File | null;
  onImageFileChange?: (file: File | null) => void;
  existingImageUrl?: string | null;
  onRemoveExistingImage?: () => void;
  isCreateForm?: boolean;
};

export default function RecipeForm({
  onSubmit,
  isSubmitting = false,
  isPreparingImageUpload = false,
  ingredients,
  onDirtyChange,
  defaultValues,
  imageFile = null,
  onImageFileChange = () => {},
  existingImageUrl = null,
  onRemoveExistingImage = () => {},
  isCreateForm = true,
}: RecipeFormProps) {
  const methods = useForm<RecipeFormValues>({
    resolver: zodResolver(recipeFormSchema),
    mode: "onChange",
    defaultValues: defaultValues ?? DEFAULT_RECIPE_FORM_VALUES,
  });

  const {
    handleSubmit,
    formState: { isDirty, isValid, isSubmitted },
  } = methods;

  useEffect(() => {
    onDirtyChange?.(isDirty);
  }, [isDirty, onDirtyChange]);

  return (
    <Container size="md">
      <FormProvider {...methods}>
        <Form onSubmit={handleSubmit(onSubmit)}>
          <Stack space="lg">
            <RecipeDetailSection
              imageFile={imageFile}
              onImageFileChange={onImageFileChange}
              existingImageUrl={existingImageUrl}
              onRemoveExistingImage={onRemoveExistingImage}
            />
            <IngredientSection ingredients={ingredients} />
            <RecipeSourceSection />
            <Inline justify="start">
            <Button disabled={isSubmitting || isPreparingImageUpload || (isSubmitted && !isValid)} type="submit" loading={isSubmitting || isPreparingImageUpload}>
              {isCreateForm ? commonStrings.forms.create : commonStrings.forms.update}
            </Button>
            </Inline>
          </Stack>
        </Form>
      </FormProvider>
    </Container>
  );
}
