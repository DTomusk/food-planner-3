import FormField from "@/components/form/FormField";
import FormSection from "@/components/form/FormSection";
import Input from "@/components/ui/Input";
import type { RecipeFormValues } from "../types";
import { useFormContext } from "react-hook-form";
import Inline from "@/components/layout/Inline";

export default function RecipeDetailSection() {
    const { register, formState: { errors } } = useFormContext<RecipeFormValues>();
    return (
        <FormSection title="Recipe details" collapsible defaultCollapsed={false}>
          <FormField htmlFor="name" label="Recipe name" error={errors.name?.message}>
            <Input type="text" placeholder="Recipe name" 
            {...register("name")} />
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