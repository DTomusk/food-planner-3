import FormField from "@/components/form/FormField";
import FormSection from "@/components/form/FormSection";
import Input from "@/components/ui/Input";
import type { RecipeFormValues } from "../types";
import { useFormContext } from "react-hook-form";
import Inline from "@/components/layout/Inline";
import FileDrop from "@/components/ui/FileDrop";
import ImageDisplay from "@/components/ui/ImageDisplay";
import Stack from "@/components/layout/Stack";
import IconButton from "@/components/ui/IconButton";
import { Trash2 } from "lucide-react";
import Label from "@/components/ui/Label";
import TextArea from "@/components/ui/TextArea";

type RecipeDetailSectionProps = {
  imageFile?: File | null;
  onImageFileChange?: (file: File | null) => void;
  existingImageUrl?: string | null;
  onRemoveExistingImage?: () => void;
};

export default function RecipeDetailSection({
  imageFile = null,
  onImageFileChange = () => {},
  existingImageUrl = null,
  onRemoveExistingImage = () => {},
}: RecipeDetailSectionProps) {
    const {
      register,
      formState: { errors },
    } = useFormContext<RecipeFormValues>();

    // Display existing image if no new image selected and there's an existing image
    const hasExistingImage = existingImageUrl !== null && !imageFile;
    const showFileDrop = !hasExistingImage;

    return (
        <FormSection title="Recipe details" collapsible defaultCollapsed={false}>
          <Label htmlFor="recipeImage">Recipe image</Label>
          {hasExistingImage && (
            <Stack space="md" className="mb-4">
              <div className="relative w-80">
                <ImageDisplay
                  imageUrl={existingImageUrl}
                  altText="Current recipe image"
                  containerClassName="w-80 h-80 rounded-md overflow-hidden"
                  imageClassName="object-contain"
                />
                <IconButton
                  aria-label="Remove existing image"
                  title="Remove image"
                  variant="danger"
                  shape="circle"
                  className="absolute right-2 top-2"
                  onClick={onRemoveExistingImage}
                >
                  <Trash2 size={16} />
                </IconButton>
              </div>
            </Stack>
          )}
          {showFileDrop && (
            <FileDrop
              value={imageFile}
              label={existingImageUrl ? "Replace recipe image" : "Upload recipe image"}
              accept={{
                "image/jpeg": [],
                "image/png": [],
              }}
              maxSize={5 * 1024 * 1024}
              onChange={onImageFileChange}
            />
          )}
          <FormField htmlFor="name" label="Recipe name" error={errors.name?.message}>
            <Input type="text" placeholder="Recipe name" 
            {...register("name")} />
          </FormField>
          <FormField htmlFor="description" label="Description" error={errors.description?.message}>
            <TextArea placeholder="Short description of the recipe (optional)" rows={3}
            {...register("description")} />
          </FormField>
          <Inline>
            <FormField htmlFor="prepMins" label="Preparation time (minutes)" error={errors.prepMins?.message}>
              <Input type="number" placeholder="Preparation time in minutes" min={1} step={1}
              {...register("prepMins", { valueAsNumber: true })} />
            </FormField>
            <FormField htmlFor="cookMins" label="Cooking time (minutes)" error={errors.cookMins?.message}>
              <Input type="number" placeholder="Cooking time in minutes" min={1} step={1}
              {...register("cookMins", { valueAsNumber: true })} />
            </FormField>
          </Inline>
          <FormField htmlFor="portions" label="Portions" error={errors.portions?.message}>
            <Input type="number" placeholder="How many portions does this make?" min={1} step={1}
            {...register("portions", { valueAsNumber: true })} />
          </FormField>
      </FormSection>  
    );
}