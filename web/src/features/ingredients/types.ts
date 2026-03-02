export type IngredientOptionModel = {
    id: string;
    name: string;
    counter: string | null | undefined;
    preferredUnit: {
        val: number;
        name: string;
        symbol: string;
    };
}