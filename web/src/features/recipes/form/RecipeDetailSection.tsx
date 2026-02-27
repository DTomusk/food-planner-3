import FormField from "@/components/form/FormField";
import FormSection from "@/components/form/FormSection";
import Input from "@/components/ui/Input";
import type { RecipeFormValues } from "../types";
import { useFormContext } from "react-hook-form";

export default function RecipeDetailSection() {
    const { register, formState: { errors } } = useFormContext<RecipeFormValues>();
    return (
        <FormSection title="Recipe details" collapsible>
          <FormField htmlFor="name" label="Recipe name" error={errors.name?.message}>
            <Input type="text" placeholder="Recipe name" 
            {...register("name", 
            { required: "Please add a name",
                minLength: { value: 3, message: "Name must be at least 3 characters" },
            })} />
          </FormField>
          <FormField htmlFor="prepMins" label="Preparation time (minutes)" error={errors.prepMins?.message}>
            <Input type="number" placeholder="Preparation time in minutes" 
            {...register("prepMins", { required: "Please add preparation time", min: { value: 1, message: "Preparation time must be at least 1 minute" } })} />
          </FormField>
          <FormField htmlFor="cookMins" label="Cooking time (minutes)" error={errors.cookMins?.message}>
            <Input type="number" placeholder="Cooking time in minutes" 
            {...register("cookMins", { required: "Please add cooking time", min: { value: 1, message: "Cooking time must be at least 1 minute" } })} />
          </FormField>
          <FormField htmlFor="portions" label="Portions" error={errors.portions?.message}>
            <Input type="number" placeholder="How many portions does this make?" 
            {...register("portions", { required: "Please add number of portions", min: { value: 1, message: "Number of portions must be at least 1" } })} />
          </FormField>
      </FormSection>  
    );
}