import type { RecipeFormValues } from "../features/recipes/types";
import { RecipeForm } from "../features/recipes/components/RecipeForm";
import PageTitle from "../components/PageTitle";
import Page from "../layout/PageWrapper";
import Alert from "../components/Alert";
import BackLink from "../components/BackLink";
import { useCreateRecipe } from "../features/recipes/hooks/useCreateRecipe";

export default function CreateRecipePage() {
    const { mutate, isPending, error } = useCreateRecipe();

    const handleSubmit = (values: RecipeFormValues) => {
        mutate(
            { input: { name: values.name } }
        );
    }

    return (
        <Page>
        <PageTitle text="Create Recipe" />
        <BackLink />
        {error && <Alert message={error.message} />}
        <RecipeForm
            onSubmit={handleSubmit}
            isSubmitting={isPending}
        />
        </Page>
    )
}