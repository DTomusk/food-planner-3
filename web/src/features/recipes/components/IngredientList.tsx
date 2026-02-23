import type { Ingredients } from "../types";

export default function IngredientList({ ingredients }: { ingredients: Ingredients[] }) {
    return (
        <ul>
            {ingredients.map((ingredient, index) => (
            <li key={index}>{ingredient.name}: {ingredient.formattedQuantity}</li>
            ))}
        </ul>
    );
}