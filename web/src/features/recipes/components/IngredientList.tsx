import { useState } from "react";
import type { Ingredient } from "../types";
import CheckListItem from "@/components/ui/CheckListItem";
import { useFormatIngredients } from "../hooks/useFormatIngredients";

export default function IngredientList({ ingredients }: { ingredients: Ingredient[] }) {
    const [checkedItems, setCheckedItems] = useState<boolean[]>(ingredients.map(() => false));
    const { formatIngredient } = useFormatIngredients();

    const handleCheckChange = (index: number) => {
        const newCheckedItems = [...checkedItems];
        newCheckedItems[index] = !newCheckedItems[index];
        setCheckedItems(newCheckedItems);
    }
    return (
        <ul className="grid gap-2">
            {ingredients.map((ingredient, index) => (
            <CheckListItem key={index} checked={checkedItems[index]} onChange={() => handleCheckChange(index)}>
                {formatIngredient(ingredient)}
            </CheckListItem>
            ))}
        </ul>
    );
}