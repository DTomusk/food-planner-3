import { Alert, BackLink, PageTitle } from "@/components";
import Container from "@/components/layout/Container";
import { useIngredients } from "@/features/ingredients/hooks/useIngredients";
import { RecipeForm, useCreateRecipe, type RecipeFormValues } from "@/features/recipes";
import { mapFormValuesToCreateRecipeInput } from "@/features/recipes/mappers/recipeFormMapper";
import { useUploadUrl } from "@/features/upload/hooks/useUploadUrl";
import { Page } from "@/layout";
import type { CreateImageUploadUrlMutation } from "@/lib/graphql.generated";
import { extractErrorMessage } from "@/lib/errors";
import { commonStrings } from "@/lib/strings";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

export default function RecipeCreatePage() {
    const { mutate, isPending, error } = useCreateRecipe();
    const { mutateAsync: createUploadUrl, isPending: isPreparingImageUpload } = useUploadUrl();
    const { data: ingredientsData } = useIngredients();
    const navigate = useNavigate();
    const [imageFile, setImageFile] = useState<File | null>(null);
    const [imageUploadPayload, setImageUploadPayload] = useState<CreateImageUploadUrlMutation["createImageUploadUrl"] | null>(null);
    const [uploadError, setUploadError] = useState<string | null>(null);

    const handleImageFileChange = async (file: File | null) => {
      setImageFile(file);
      setImageUploadPayload(null);
      setUploadError(null);

      if (!file) {
        return;
      }

      try {
        const data = await createUploadUrl({
          input: {
            fileName: file.name,
            fileType: file.type,
            fileSize: file.size,
          },
        });
        setImageUploadPayload(data.createImageUploadUrl);
      } catch (err) {
        setUploadError(extractErrorMessage(err));
      }
    };

    const handleSubmit = (values: RecipeFormValues) => {
      if (imageFile && !imageUploadPayload) {
        setUploadError("Unable to prepare image upload. Please re-select your image and try again.");
        return;
      }

      mutate(
          { input: mapFormValuesToCreateRecipeInput(values) },
          {
            onSuccess: (data) => {
              navigate(`/recipes/${data.createRecipe.id}`, {
                state: {
                  successMessage: "Recipe created successfully!",
                  imageUploadId: imageUploadPayload?.uploadId,
                  imageFileUrl: imageUploadPayload?.fileUrl,
                }
              });
            }
          }
      );
    };

    return (
        <Page toolbarLeft={<BackLink />} >
            <PageTitle text={commonStrings.recipe.create} />
          {uploadError && <Container><Alert message={uploadError} closable /></Container>}
            {error && <Container><Alert message={extractErrorMessage(error)} closable /></Container>}
          <RecipeForm
            onSubmit={handleSubmit}
            isSubmitting={isPending}
            isPreparingImageUpload={isPreparingImageUpload}
            ingredients={ingredientsData || []}
            imageFile={imageFile}
            onImageFileChange={handleImageFileChange}
          />
        </Page>
    );
}