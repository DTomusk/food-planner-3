import type { Recipe } from "../types";

interface RecipeCardProps {
    recipe: Recipe;
    onClick?: () => void;
}

export default function RecipeCard({ recipe, onClick }: RecipeCardProps) {
    return (
        <div className="border p-4 rounded shadow" onClick={onClick}>
            <h3 className="text-xl font-semibold">{recipe.name}</h3>
        </div>
    );
}