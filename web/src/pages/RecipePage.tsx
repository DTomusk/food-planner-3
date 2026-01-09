import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { recipeQuery } from "../features/recipes/queries";
import PageTitle from "../components/PageTitle";
import RecipeCard from "../features/recipes/components/RecipeCard";
import Page from "../layout/PageWrapper";
import BackLink from "../components/BackLink";
import Spinner from "../components/Spinner";
import ErrorAlert from "../components/ErrorAlert";

export default function RecipePage() {
    const { id } = useParams<{ id: string }>();
    const { data, isLoading, error } = useQuery(recipeQuery(id!));

    return (
        <Page>
            <PageTitle text="Recipe Page" />
            <BackLink />
            {!id && <ErrorAlert message="No recipe ID provided." />}
            {isLoading && <Spinner />}
            {error && <ErrorAlert message={(error as Error).message} />}
            {data?.recipe ? (
                <RecipeCard recipe={data.recipe} />
            ) : (
                !isLoading && id && !error && <ErrorAlert message="Recipe not found." />
            )}
        </Page>
    )
}