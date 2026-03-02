export function useFormatIngredientOption(ingredientName: string, counter: string) {
    if (!counter) return ingredientName;
    return `${ingredientName} (${counter})`;
}