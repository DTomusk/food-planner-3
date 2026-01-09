import { useMutation } from "@tanstack/react-query";
import type { RecipeFormValues } from "../features/recipes/types";
import { createRecipe } from "../features/recipes/mutations";
import { RecipeForm } from "../features/recipes/components/RecipeForm";
import PageTitle from "../components/PageTitle";
import { queryClient } from "../app/queryClient";
import Page from "../layout/PageWrapper";

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
        <Page>
        <PageTitle text="Create Recipe" />
        <RecipeForm
            onSubmit={(values) => mutation.mutate(values)}
            isSubmitting={mutation.isPending}
        />
        </Page>
    )
}