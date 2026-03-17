import Stack from "@/components/layout/Stack";
import { RecipeSourceType } from "../types";

const SOURCE_TYPE_OPTIONS: Array<{
    value: RecipeSourceType;
    label: string;
}> = [
    {
        value: RecipeSourceType.None,
        label: "No source - don't include any source information",
    },
    {
        value: RecipeSourceType.Website,
        label: "Website - add a link to the recipe",
    },
    {
        value: RecipeSourceType.Cookbook,
        label: "Cookbook - add the name of the cookbook and the page number of the recipe",
    },
    {
        value: RecipeSourceType.Original,
        label: "Original recipe - add your own instructions for the recipe",
    },
];

interface RecipeSourceTypeSelectorProps {
    name: string;
    value: RecipeSourceType | undefined;
    onChange: (value: RecipeSourceType) => void;
};

export default function RecipeSourceTypeSelector({ name, value, onChange }: RecipeSourceTypeSelectorProps) {
    return (
        <Stack space="sm">
            {SOURCE_TYPE_OPTIONS.map((option) => (
                <label key={option.value} className="flex items-center gap-2">
                    <input
                        type="radio"
                        name={name}
                        value={option.value}
                        checked={value === option.value}
                        onChange={() => onChange(option.value)}
                    />
                    {option.label}
                </label>
            ))}
        </Stack>
    );
}
