import type { Ingredient } from "../types";

// TODO: this can be a work in progress
export function useFormatIngredients() {
    const formatIngredient = (ingredient : Ingredient) => {
        const lowercaseName = ingredient.name.toLowerCase();   
        const formattedQuantity = ingredient.quantity + ingredient.unitSymbol;
        if (ingredient.counter) return `${formattedQuantity} ${ingredient.counter} of ${lowercaseName}`;
        return `${formattedQuantity} ${lowercaseName}`;
    }

    return { formatIngredient, };
}