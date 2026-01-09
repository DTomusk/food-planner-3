import type { Recipe } from "../types";

export default function RecipeCard({ recipe }: { recipe: Recipe }) {
    return (
        <div className="border p-4 rounded shadow">
            <h3 className="text-xl font-semibold">{recipe.name}</h3>
        </div>
    );
}