import { Alert, PageTitle, Spinner } from "@/components";
import Container from "@/components/layout/Container";
import { useIngredients } from "@/features/ingredients/hooks/useIngredients";
import { RecipeForm, useRecipe, type RecipeFormValues } from "@/features/recipes";
import { useUpdateRecipe } from "@/features/recipes/hooks/useUpdateRecipe";
import { mapFormValuesToUpdateRecipeInput } from "@/features/recipes/mappers/recipeFormMapper";
import { useUploadFileToSignedUrl } from "@/features/upload/hooks/useUploadFileToSignedUrl";
import { useUploadUrl } from "@/features/upload/hooks/useUploadUrl";
import { Page } from "@/layout";
import { commonStrings } from "@/lib";
import { extractErrorMessage } from "@/lib/errors";
import type { CreateImageUploadUrlMutation } from "@/lib/graphql.generated";
import imageCompression from "browser-image-compression";
import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";

export default function RecipeUpdatePage() {
    const { id } = useParams<{ id: string }>();
    const { data, isLoading, error } = useRecipe(id ?? "");
    const { mutate, isPending, error: mutateError } = useUpdateRecipe();
    const { data: ingredientsData } = useIngredients();
    const { mutateAsync: createUploadUrl, isPending: isPreparingImageUpload } = useUploadUrl();
    const { mutateAsync: uploadFileToSignedUrl, isPending: isUploadingImage } = useUploadFileToSignedUrl();
    const navigate = useNavigate();
    
    const [imageFile, setImageFile] = useState<File | null>(null);
    const [imageUploadPayload, setImageUploadPayload] = useState<CreateImageUploadUrlMutation["createImageUploadUrl"] | null>(null);
    const [uploadError, setUploadError] = useState<string | null>(null);
    const [isCompressingImage, setIsCompressingImage] = useState(false);
    const [hasRemovedExistingImage, setHasRemovedExistingImage] = useState(false);

    const existingImageUrl = data?.formValues?.imgSrc ?? null;

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

    const handleRemoveExistingImage = () => {
        setHasRemovedExistingImage(true);
        setImageFile(null);
        setImageUploadPayload(null);
        setUploadError(null);
    };

    const handleSubmit = (values: RecipeFormValues) => {
        if (!id) {
            return;
        }

        if (imageFile && !imageUploadPayload) {
            setUploadError("Image is not uploaded yet. Please re-select your image and try again.");
            return;
        }

        const imageChangedOrRemoved = hasRemovedExistingImage;

        const { input } = mapFormValuesToUpdateRecipeInput(values, {
            imgUploadId: imageUploadPayload?.uploadId,
            removeImage: imageChangedOrRemoved,
        });

        mutate(
            {
                input: {
                    id,
                    details: input,
                    removeImage: imageChangedOrRemoved ? true : undefined,
                },
            },
            {
                onSuccess: (data) => {
                    navigate(`/recipes/${data.updateRecipe.id}`, {
                        state: { successMessage: "Recipe updated successfully!" }
                    });
                }
            }
        );
    };

    return (
        <Page>
            <PageTitle text={commonStrings.recipe.update} />
            {!id && <Container><Alert message="No recipe ID provided." /></Container>}
            {isLoading && <Container><Spinner /></Container>}
            {error && <Container><Alert message={extractErrorMessage(error)} closable /></Container>}
            {mutateError && <Container><Alert message={extractErrorMessage(error)} closable /></Container>}
            {uploadError && <Container><Alert message={uploadError} closable /></Container>}
            {id && data?.formValues && !isLoading && (
                <RecipeForm
                    key={id}
                    onSubmit={handleSubmit}
                    isSubmitting={isPending}
                    isPreparingImageUpload={isCompressingImage || isPreparingImageUpload || isUploadingImage}
                    ingredients={ingredientsData || []}
                    defaultValues={data.formValues}
                    imageFile={imageFile}
                    onImageFileChange={handleImageFileChange}
                    existingImageUrl={hasRemovedExistingImage ? null : existingImageUrl}
                    onRemoveExistingImage={handleRemoveExistingImage}
                    isCreateForm={false}
                />
            )}
        </Page>
    );
}