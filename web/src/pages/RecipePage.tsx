import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { recipeQuery } from "../features/recipes/queries";
import Title from "../components/Title";
import RecipeCard from "../features/recipes/components/RecipeCard";

export default function RecipePage() {
    const { id } = useParams<{ id: string }>();

    if (!id) {
        return <p>No recipe ID provided.</p>;
    }

    const { data, isLoading, error } = useQuery(recipeQuery(id));

    if (isLoading) {
        return <p>Loading...</p>;
    }

    if (error) {
        return <p>Error: {(error as Error).message}</p>;
    }

    if (!data.recipe) {
        return <p>Recipe not found.</p>;
    }

    return (
        <>
            <Title text="Recipe Page" />
            {data && (
                <RecipeCard recipe={data.recipe} />
            )}
        </>
    )
}