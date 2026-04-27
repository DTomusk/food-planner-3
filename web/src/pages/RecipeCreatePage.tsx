import { Alert, PageTitle } from "@/components";
import Container from "@/components/layout/Container";
import { useIngredients } from "@/features/ingredients/hooks/useIngredients";
import { RecipeForm, useCreateRecipe, type RecipeFormValues } from "@/features/recipes";
import { mapFormValuesToCreateRecipeInput } from "@/features/recipes/mappers/recipeFormMapper";
import { useUploadFileToSignedUrl } from "@/features/upload/hooks/useUploadFileToSignedUrl";
import { useUploadUrl } from "@/features/upload/hooks/useUploadUrl";
import { Page } from "@/layout";
import type { CreateImageUploadUrlMutation } from "@/lib/graphql.generated";
import { extractErrorMessage } from "@/lib/errors";
import { useUnsavedChanges } from "@/lib/hooks/useUnsavedChanges";
import { commonStrings } from "@/lib/strings";
import imageCompression from "browser-image-compression";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

export default function RecipeCreatePage() {
  const { mutateAsync: createRecipe, isPending, error } = useCreateRecipe();
  const { mutateAsync: createUploadUrl, isPending: isPreparingImageUpload } = useUploadUrl();
  const { mutateAsync: uploadFileToSignedUrl, isPending: isUploadingImage } = useUploadFileToSignedUrl();
  const { data: ingredientsData } = useIngredients();
  const navigate = useNavigate();
  const [imageFile, setImageFile] = useState<File | null>(null);
  const [imageUploadPayload, setImageUploadPayload] = useState<CreateImageUploadUrlMutation["createImageUploadUrl"] | null>(null);
  const [uploadError, setUploadError] = useState<string | null>(null);
  const [isCompressingImage, setIsCompressingImage] = useState(false);
  const [isFormDirty, setIsFormDirty] = useState(false);
  const shouldWarnUnsavedChanges = !isPending && (isFormDirty || Boolean(imageFile));

  // Image handling lives outside RHF, so include selected file state in the unsaved check.
  useUnsavedChanges(shouldWarnUnsavedChanges, {
    blockRouteChange: true,
  });

  const handleImageFileChange = async (file: File | null) => {
    setImageFile(file);
    setImageUploadPayload(null);
    setUploadError(null);

    if (!file) {
      setIsCompressingImage(false);
      return;
    }

    try {
      setIsCompressingImage(true);

      // Compress image before uploading
      const compressedFile = await imageCompression(file, {
        maxSizeMB: 1,
        maxWidthOrHeight: 2048,
        useWebWorker: true,
      });

      setIsCompressingImage(false);

      const data = await createUploadUrl({
        input: {
          fileName: file.name,
          fileType: compressedFile.type,
          fileSize: compressedFile.size,
        },
      });

      const signedUpload = data.createImageUploadUrl;

      await uploadFileToSignedUrl({
        uploadUrl: signedUpload.uploadUrl,
        file: compressedFile,
        fileType: compressedFile.type,
      });

      setImageUploadPayload(signedUpload);
    } catch (err) {
      setIsCompressingImage(false);
      setUploadError(extractErrorMessage(err));
    }
  };

  const handleSubmit = async (values: RecipeFormValues, options: { publish: boolean }) => {
    if (imageFile && !imageUploadPayload) {
      setUploadError("Image is not uploaded yet. Please re-select your image and try again.");
      return;
    }

    try {
      const data = await createRecipe({
        input: mapFormValuesToCreateRecipeInput(values, {
          imgUploadId: imageUploadPayload?.uploadId,
          publish: options.publish,
        }),
      });

      navigate(`/recipes/${data.createRecipe.id}`, {
        state: {
          successMessage: "Recipe created successfully!",
          imageUploadId: imageUploadPayload?.uploadId,
          imageFileUrl: imageUploadPayload?.fileUrl,
        }
      });
    } catch {
      // Create recipe errors are handled by useCreateRecipe state.
    }
  };

  return (
    <Page >
      <PageTitle text={commonStrings.recipe.create} />
      {uploadError && <Container><Alert message={uploadError} closable /></Container>}
      {error && <Container><Alert message={extractErrorMessage(error)} closable /></Container>}
      <RecipeForm
        onSubmit={handleSubmit}
        isSubmitting={isPending}
        isPreparingImageUpload={isCompressingImage || isPreparingImageUpload || isUploadingImage}
        ingredients={ingredientsData || []}
        onDirtyChange={setIsFormDirty}
        imageFile={imageFile}
        onImageFileChange={handleImageFileChange}
      />
    </Page>
  );
}