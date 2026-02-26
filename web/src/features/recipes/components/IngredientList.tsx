import { useState } from "react";
import type { Ingredients } from "../types";
import CheckListItem from "@/components/ui/CheckListItem";

export default function IngredientList({ ingredients }: { ingredients: Ingredients[] }) {
    const [checkedItems, setCheckedItems] = useState<boolean[]>(ingredients.map(() => false));

    const handleCheckChange = (index: number) => {
        const newCheckedItems = [...checkedItems];
        newCheckedItems[index] = !newCheckedItems[index];
        setCheckedItems(newCheckedItems);
    }
    return (
        <ul className="grid gap-2">
            {ingredients.map((ingredient, index) => (
            <CheckListItem key={index} checked={checkedItems[index]} onChange={() => handleCheckChange(index)}>
                {ingredient.name}: {ingredient.formattedQuantity} 
            </CheckListItem>
            ))}
        </ul>
    );
}