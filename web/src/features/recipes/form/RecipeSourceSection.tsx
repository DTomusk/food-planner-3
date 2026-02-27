import { Controller, useFormContext, useWatch } from "react-hook-form";
import type { RecipeFormValues } from "../types";
import FormSection from "@/components/form/FormSection";
import FormField from "@/components/form/FormField";
import Input from "@/components/ui/Input";
import Inline from "@/components/layout/Inline";
import Stack from "@/components/layout/Stack";
import SectionHelpText from "@/components/form/SectionHelpText";
import TextArea from "@/components/ui/TextArea";

export default function RecipeSourceSection() {
    const { control, register, formState: { errors } } = useFormContext<RecipeFormValues>();
    const sourceType = useWatch({
        control,
        name: "sourceType",
    });
    
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
                    <Stack space="sm">
                    <label className="flex items-center gap-2">
                        <input
                        type="radio"
                        checked={field.value === 0}
                        onChange={() => field.onChange(0)}
                        />
                        No source - don't include any source information
                    </label>  
                    <label className="flex items-center gap-2">
                        <input
                        type="radio"
                        checked={field.value === 1}
                        onChange={() => field.onChange(1)}
                        />
                        Website - add a link to the recipe
                    </label>

                    <label className="flex items-center gap-2">
                        <input
                        type="radio"
                        checked={field.value === 2}
                        onChange={() => field.onChange(2)}
                        />
                        Cookbook - add the name of the cookbook and the page number of the recipe
                    </label>

                    <label className="flex items-center gap-2">
                        <input
                        type="radio"
                        checked={field.value === 3}
                        onChange={() => field.onChange(3)}
                        />
                        Original recipe - add your own instructions for the recipe
                    </label>
                    </Stack>
                )}
                />
            {sourceType === 1 && (
                <FormField htmlFor="url" label="URL" error={errors.url?.message}>
                    <Input type="text" placeholder="URL of the recipe" {...register("url")} />
                </FormField>
            )}
            {sourceType === 2 && (
                <Inline>
                    <FormField htmlFor="cookbook" label="Cookbook" error={errors.bookTitle?.message}>
                        <Input type="text" placeholder="Name of the cookbook" {...register("bookTitle")} />
                    </FormField>
                    <FormField htmlFor="page" label="Page number" error={errors.bookPage?.message}>
                        <Input type="number" placeholder="Page number" {...register("bookPage")} />
                    </FormField>
                </Inline>
            )}
            {sourceType === 3 && (
                <FormField htmlFor="instructions" label="Instructions" error={errors.instructions?.message}>
                <TextArea
                placeholder="Add the instructions for the recipe"
                {...register("instructions")}
                />
                </FormField>
            )}
        </FormSection>
    );
}