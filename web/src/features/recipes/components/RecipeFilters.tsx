import { Stack, Text } from "@/components";
import CollapsibleHeader from "@/components/ui/CollapsibleHeader";
import Label from "@/components/ui/Label";
import Select from "@/components/ui/Select";
import { useState } from "react";

type RecipeDietLevelFilter = "all" | "0" | "1";

type RecipeFiltersProps = {
    dietLevel: RecipeDietLevelFilter;
    onDietLevelChange: (dietLevel: RecipeDietLevelFilter) => void;
    glutenFreeOnly: boolean;
    onGlutenFreeOnlyChange: (glutenFreeOnly: boolean) => void;
    mobile?: boolean;
};

const dietLevelOptions: Array<{ value: RecipeDietLevelFilter; label: string; }> = [
    {
        value: "all",
        label: "Any",
    },
    {
        value: "0",
        label: "Vegan",
    },
    {
        value: "1",
        label: "Vegetarian",
    },
];

export default function RecipeFilters({ dietLevel, onDietLevelChange, glutenFreeOnly, onGlutenFreeOnlyChange, mobile = false }: RecipeFiltersProps) {
    const [collapsed, setCollapsed] = useState(false);
    return (
        <>
        {mobile && 
        (<Stack space="sm">
            <Text>Filters</Text>
            <Label htmlFor="dietLevel" >
                Diet
            </Label>
            <Select id="dietLevel" value={dietLevel} onChange={(e) => onDietLevelChange(e.target.value as RecipeDietLevelFilter)}>
                {dietLevelOptions.map((option) => (
                    <option key={option.value} value={option.value}>
                        {option.label}
                    </option>
                ))}
            </Select>
            <label htmlFor="glutenFreeOnlyMobile" className="flex items-center gap-2 text-sm text-gray-900">
                <input
                    id="glutenFreeOnlyMobile"
                    type="checkbox"
                    checked={glutenFreeOnly}
                    onChange={(e) => onGlutenFreeOnlyChange(e.target.checked)}
                />
                Gluten-free only
            </label>
        </Stack>)
        }
        {!mobile && <aside className="w-full rounded-md border border-black bg-white p-3 md:w-56 md:shrink-0">
            <fieldset>
                <CollapsibleHeader title="Filters" collapsible collapsed={collapsed} onToggle={setCollapsed} />
                
                {!collapsed && (
                    <>
                    <Stack space="md">
                    <Label htmlFor="dietLevel" >
                        Diet
                    </Label>
                    <Select id="dietLevel" value={dietLevel} onChange={(e) => onDietLevelChange(e.target.value as RecipeDietLevelFilter)}>
                        {dietLevelOptions.map((option) => (
                            <option key={option.value} value={option.value}>
                                {option.label}
                            </option>
                        ))}
                    </Select>

                    <label htmlFor="glutenFreeOnlyDesktop" className="flex items-center gap-2 text-sm text-gray-900">
                        <input
                            id="glutenFreeOnlyDesktop"
                            type="checkbox"
                            checked={glutenFreeOnly}
                            onChange={(e) => onGlutenFreeOnlyChange(e.target.checked)}
                        />
                        Gluten-free only
                    </label>
                    </Stack>
                    </>
                )}
            </fieldset>
        </aside>
        }
        </>
    );
}