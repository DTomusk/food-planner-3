type RecipeDietLevelFilter = "all" | "0" | "1";

type RecipeFiltersProps = {
    dietLevel: RecipeDietLevelFilter;
    onDietLevelChange: (dietLevel: RecipeDietLevelFilter) => void;
};

const dietLevelOptions: Array<{ value: RecipeDietLevelFilter; label: string; description: string }> = [
    {
        value: "all",
        label: "All recipes",
        description: "No diet restriction",
    },
    {
        value: "0",
        label: "Vegan",
        description: "Only vegan recipes",
    },
    {
        value: "1",
        label: "Vegetarian",
        description: "Vegan and vegetarian recipes",
    },
];

export default function RecipeFilters({ dietLevel, onDietLevelChange }: RecipeFiltersProps) {
    return (
        <aside className="w-full rounded-md border border-black bg-white p-4 md:w-56 md:shrink-0">
            <fieldset>
                <legend className="mb-3 text-sm font-semibold text-gray-900">Filters</legend>
                <div className="space-y-2">
                    {dietLevelOptions.map((option) => (
                        <label
                            key={option.value}
                            className="flex cursor-pointer items-start gap-2 rounded px-1 py-1 hover:bg-gray-50"
                        >
                            <input
                                type="radio"
                                name="diet-level"
                                value={option.value}
                                checked={dietLevel === option.value}
                                onChange={() => onDietLevelChange(option.value)}
                                className="mt-0.5"
                            />
                            <span>
                                <span className="block text-sm text-gray-900">{option.label}</span>
                                <span className="block text-xs text-gray-500">{option.description}</span>
                            </span>
                        </label>
                    ))}
                </div>
            </fieldset>
        </aside>
    );
}