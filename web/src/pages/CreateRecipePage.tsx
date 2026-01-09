import { useMutation } from "@tanstack/react-query";
import type { RecipeFormValues } from "../features/recipes/types";
import { createRecipe } from "../features/recipes/mutations";
import { RecipeForm } from "../features/recipes/components/RecipeForm";
import PageTitle from "../components/PageTitle";
import { queryClient } from "../app/queryClient";
import Page from "../layout/PageWrapper";
import { useState } from "react";
import ErrorAlert from "../components/ErrorAlert";

export default function CreateRecipePage() {
    const [error, setError] = useState<string | null>(null);
    const mutation = useMutation({
        mutationFn: (values: RecipeFormValues) => createRecipe(values),
        onSuccess: () => {
            alert("Recipe created successfully!");
            queryClient.invalidateQueries({ queryKey: ['recipes'] });
            setError(null);
        },
        onError: (error) => {
            alert(`Error creating recipe: ${(error as Error).message}`);
            setError((error as Error).message);
        },
    })

    return (
        <Page>
        <PageTitle text="Create Recipe" />
        {error && <ErrorAlert message={error} />}
        <RecipeForm
            onSubmit={(values) => mutation.mutate(values)}
            isSubmitting={mutation.isPending}
        />
        </Page>
    )
}