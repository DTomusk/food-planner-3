import type { Ingredient } from "../types";

// TODO: this can be a work in progress
export function useFormatIngredients() {
    const formatIngredient = (ingredient : Ingredient) => {
        const lowercaseName = ingredient.name.toLowerCase();   

        // For now simply quantity + unit + name, e.g. 200g flour
        if (ingredient.unitSymbol.length > 0) {
            return `${ingredient.quantity}${ingredient.unitSymbol} ${lowercaseName}`;
        }

        // The following are unitless, so can be pluralised
        const isPlural = ingredient.quantity != 1;

        if (ingredient.counter) return `${ingredient.quantity} ${isPlural ? ingredient.counterPlural : ingredient.counter} of ${lowercaseName}`;

        return `${ingredient.quantity} ${isPlural ? ingredient.plural?.toLowerCase() : lowercaseName}`;
    }

    return { formatIngredient, };
}