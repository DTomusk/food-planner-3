import { useParams } from "react-router-dom";
import PageTitle from "../components/PageTitle";
import RecipeCard from "../features/recipes/components/RecipeCard";
import Page from "../layout/PageWrapper";
import BackLink from "../components/BackLink";
import Spinner from "../components/Spinner";
import Alert from "../components/Alert";
import { useRecipe } from "../features/recipes/hooks/useRecipe";

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