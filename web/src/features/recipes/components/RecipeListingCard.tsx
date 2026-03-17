import Inline from "@/components/layout/Inline";
import type { RecipeSummary } from "../types";

interface RecipeListingCardProps {
    recipe: RecipeSummary;
    onClick?: () => void;
    actions?: React.ReactNode;
}

export default function RecipeListingCard({ recipe, onClick, actions }: RecipeListingCardProps) {
    return (
        <div className="border bg-white p-4 rounded shadow cursor-pointer hover:scale-105 transition-transform" onClick={onClick}>
            <Inline justify="between" align="center" wrap>
                <h3 className="text-xl font-semibold">{recipe.name}</h3>
                {actions && <Inline>{actions}</Inline>}
            </Inline>
        </div>
    );
}