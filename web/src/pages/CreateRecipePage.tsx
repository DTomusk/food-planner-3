import { useMutation } from "@tanstack/react-query";
import type { RecipeFormValues } from "../features/recipes/types";
import { createRecipe } from "../features/recipes/mutations";
import { RecipeForm } from "../features/recipes/components/RecipeForm";
import Title from "../components/Title";
import { queryClient } from "../app/queryClient";

export default function CreateRecipePage() {
    const mutation = useMutation({
        mutationFn: (values: RecipeFormValues) => createRecipe(values),
        onSuccess: () => {
            alert("Recipe created successfully!");
            queryClient.invalidateQueries({ queryKey: ['recipes'] });
        },
        onError: (error) => {
            alert(`Error creating recipe: ${(error as Error).message}`);
        },
    })

    return (
        <>
        <Title text="Create Recipe" />
        <RecipeForm
            onSubmit={(values) => mutation.mutate(values)}
            isSubmitting={mutation.isPending}
        />
        </>
    )
}