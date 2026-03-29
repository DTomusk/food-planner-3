import Inline from "@/components/layout/Inline";
import type { RecipeSummary } from "../types";
import ImageDisplay from "@/components/ui/ImageDisplay";

interface RecipeListingCardProps {
    recipe: RecipeSummary;
    onClick?: () => void;
    actions?: React.ReactNode;
}

export default function RecipeListingCard({ recipe, onClick, actions }: RecipeListingCardProps) {
    return (
        <div className="border bg-white rounded shadow cursor-pointer hover:scale-105 transition-transform" onClick={onClick}>
            <Inline>
                <ImageDisplay imageUrl={recipe.imageUrl} altText={recipe.name} />
                <Inline justify="between" align="center" wrap>
                    <h3 className="text-xl font-semibold">{recipe.name}</h3>
                    {actions && <Inline>{actions}</Inline>}
                </Inline>
            </Inline>
        </div>
    );
}