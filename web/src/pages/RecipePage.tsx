import { useParams } from "react-router-dom";
import {Alert, BackLink, PageTitle, Spinner } from "../components";
import { RecipeCard, useRecipe } from "../features/recipes";
import { Page } from "../layout";

export default function RecipePage() {
    const { id } = useParams<{ id: string }>();
    const { data, isLoading, error } = useRecipe(id!);

    return (
        <Page>
            <PageTitle text="Recipe Page" />
            <BackLink />
            {!id && <Alert message="No recipe ID provided." />}
            {isLoading && <Spinner />}
            {error && <Alert message={(error as Error).message} />}
            {data ? (
                <RecipeCard recipe={data} />
            ) : (
                !isLoading && id && !error && <Alert message="Recipe not found." />
            )}
        </Page>
    )
}