import type { RecipeSummary } from "../types";

interface RecipeListingCardProps {
    recipe: RecipeSummary;
    onClick?: () => void;
}

export default function RecipeListingCard({ recipe, onClick }: RecipeListingCardProps) {
    return (
        <div className="border p-4 rounded shadow cursor-pointer hover:scale-105 transition-transform" onClick={onClick}>
            <h3 className="text-xl font-semibold">{recipe.name}</h3>
        </div>
    );
}