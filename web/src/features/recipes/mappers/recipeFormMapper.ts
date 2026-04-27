import type { CreateRecipeInput, GetRecipeQuery } from "@/lib/graphql.generated";
import { RecipeSourceType, type RecipeFormValues } from "../types";

export const DEFAULT_RECIPE_FORM_VALUES: RecipeFormValues = {
  name: "",
  prepMins: 0,
  cookMins: 0,
  portions: 0,
  ingredientUsages: [{ ingredientId: "", quantity: 0, unit: 1 }],
  sourceType: RecipeSourceType.None,
  url: "",
  bookTitle: "",
  bookPage: undefined,
  instructions: "",
};

export function toRecipeSourceType(sourceType: number): RecipeSourceType {
  switch (sourceType) {
    case 1:
      return RecipeSourceType.Website;
    case 2:
      return RecipeSourceType.Cookbook;
    case 3:
      return RecipeSourceType.Original;
    default:
      return RecipeSourceType.None;
  }
}

export function toRecipeSourceTypeValue(sourceType: RecipeSourceType): number {
  switch (sourceType) {
    case RecipeSourceType.Website:
      return 1;
    case RecipeSourceType.Cookbook:
      return 2;
    case RecipeSourceType.Original:
      return 3;
    default:
      return 0;
  }
}

type MapCreateRecipeInputOptions = {
  imgUploadId?: string | null;
  publish: boolean;
};

type MapUpdateRecipeInputOptions = {
  imgUploadId?: string | null;
  removeImage?: boolean;
  publish: boolean;
};

export function mapFormValuesToCreateRecipeInput(
  values: RecipeFormValues,
  options: MapCreateRecipeInputOptions,
): CreateRecipeInput {
  return {
    name: values.name,
    description: values.description?.trim() ? values.description.trim() : undefined,
    ingredientUsages: values.ingredientUsages.map((usage) => ({
      ingredientID: usage.ingredientId,
      quantity: usage.quantity,
      unit: usage.unit,
    })),
    prepMins: values.prepMins,
    cookMins: values.cookMins,
    portions: values.portions,
    recipeSource: {
      type: toRecipeSourceTypeValue(values.sourceType),
      url: values.sourceType === RecipeSourceType.Website ? values.url : undefined,
      bookTitle: values.sourceType === RecipeSourceType.Cookbook ? values.bookTitle : undefined,
      bookPage: values.sourceType === RecipeSourceType.Cookbook ? values.bookPage : undefined,
      instructions: values.sourceType === RecipeSourceType.Original ? values.instructions : undefined,
    },
    imgUploadId: options?.imgUploadId ?? undefined,
    publish: options.publish,
  };
}

export function mapFormValuesToUpdateRecipeInput(
  values: RecipeFormValues,
  options: MapUpdateRecipeInputOptions,
): { input: CreateRecipeInput; removeImage?: boolean } {
  return {
    input: mapFormValuesToCreateRecipeInput(values, {
      imgUploadId: options?.imgUploadId,
      publish: options.publish,
    }),
    removeImage: options?.removeImage,
  };
}

export function mapRecipeToFormValues(
  recipe: NonNullable<GetRecipeQuery["recipe"]>
): RecipeFormValues {
  const ingredientUsages = recipe.currentVersion.ingredientUsages.map((usage) => ({
    ingredientId: usage.ingredient.id,
    quantity: usage.quantity,
    unit: usage.unit.val,
  }));

  return {
    name: recipe.currentVersion.name,
    description: recipe.currentVersion.description ?? "",
    prepMins: recipe.currentVersion.prepMins,
    cookMins: recipe.currentVersion.cookMins,
    portions: recipe.currentVersion.portions,
    ingredientUsages:
      ingredientUsages.length > 0
        ? ingredientUsages
        : [{ ...DEFAULT_RECIPE_FORM_VALUES.ingredientUsages[0] }],
    sourceType: toRecipeSourceType(recipe.currentVersion.source.type),
    url: recipe.currentVersion.source.url ?? "",
    bookTitle: recipe.currentVersion.source.bookTitle ?? "",
    bookPage: recipe.currentVersion.source.bookPage ?? undefined,
    instructions: recipe.currentVersion.source.instructions ?? "",
    imgSrc: recipe.currentVersion.imgSrc ?? undefined,
  };
}