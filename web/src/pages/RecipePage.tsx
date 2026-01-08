import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { recipeQuery } from "../features/recipes/queries";

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
            <h1 className="text-3xl font-bold underline">
                Recipe Page
            </h1>
            {data && (
                <div>
                    <h2 className="text-2xl font-bold">{data.recipe.name}</h2>
                    <p>ID: {data.recipe.id}</p>
                </div>
            )}
        </>
    )
}