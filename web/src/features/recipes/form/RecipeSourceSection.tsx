import { Controller, useFormContext, useWatch } from "react-hook-form";
import { RecipeSourceType, type RecipeFormValues } from "../types";
import FormSection from "@/components/form/FormSection";
import FormField from "@/components/form/FormField";
import Input from "@/components/ui/Input";
import Inline from "@/components/layout/Inline";
import SectionHelpText from "@/components/form/SectionHelpText";
import TextArea from "@/components/ui/TextArea";
import RecipeSourceTypeSelector from "./RecipeSourceTypeSelector";

export default function RecipeSourceSection() {
    const { control, register, resetField, formState: { errors } } = useFormContext<RecipeFormValues>();
    const sourceType = useWatch({
        control,
        name: "sourceType",
    });

    function handleSourceTypeChange(onChange: (value: RecipeSourceType) => void, next: RecipeSourceType) {
        if (sourceType && sourceType !== next) {
            switch (sourceType) {
                case RecipeSourceType.Website:
                    resetField("url", { defaultValue: "" });
                    break;
                case RecipeSourceType.Cookbook:
                    resetField("bookTitle", { defaultValue: "" });
                    resetField("bookPage", { defaultValue: undefined });
                    break;
                case RecipeSourceType.Original:
                    resetField("instructions", { defaultValue: "" });
                    break;
                default:
                    break;
            }
        }
        onChange(next);
    }

    return (
        <FormSection title="Source/instructions (optional)" collapsible defaultCollapsed >
            <SectionHelpText>
                Please choose where the recipe comes from. If it's your own, you can add the instructions here. Otherwise, you can reference the site or book the recipe comes from.
            </SectionHelpText>
            <Controller
                name="sourceType"
                control={control}
                rules={{ required: "Please select a source type" }}
                render={({ field }) => (
                    <RecipeSourceTypeSelector
                        name={field.name}
                        value={field.value}
                        onChange={(next) => handleSourceTypeChange(field.onChange, next)}
                    />
                )}
                />
            {sourceType === RecipeSourceType.Website && (
                <FormField htmlFor="url" label="URL" error={errors.url?.message}>
                    <Input type="text" placeholder="URL of the recipe" {...register("url")} />
                </FormField>
            )}
            {sourceType === RecipeSourceType.Cookbook && (
                <Inline>
                    <FormField htmlFor="cookbook" label="Cookbook" error={errors.bookTitle?.message}>
                        <Input type="text" placeholder="Name of the cookbook" {...register("bookTitle")} />
                    </FormField>
                    <FormField htmlFor="page" label="Page number" error={errors.bookPage?.message}>
                        <Input type="number" placeholder="Page number" {...register("bookPage", { setValueAs: (v) => v === "" ? undefined : parseInt(v, 10) })} />
                    </FormField>
                </Inline>
            )}
            {sourceType === RecipeSourceType.Original && (
                <FormField htmlFor="instructions" label="Instructions" error={errors.instructions?.message}>
                <TextArea
                placeholder="Add the instructions for the recipe"
                rows={10}
                {...register("instructions")}
                />
                </FormField>
            )}
        </FormSection>
    );
}