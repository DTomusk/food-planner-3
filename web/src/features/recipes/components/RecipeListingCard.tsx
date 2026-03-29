import type { RecipeSummary } from "../types";
import ImageDisplay from "@/components/ui/ImageDisplay";

interface RecipeListingCardProps {
    recipe: RecipeSummary;
    onClick?: () => void;
    actions?: React.ReactNode;
}

export default function RecipeListingCard({ recipe, onClick, actions }: RecipeListingCardProps) {
    return (
        <div
            className="cursor-pointer overflow-hidden rounded border bg-white shadow transition-transform hover:scale-[1.02]"
            onClick={onClick}
        >
            <div className="flex flex-col md:flex-row">
                <ImageDisplay
                    imageUrl={recipe.imageUrl}
                    altText={recipe.name}
                    containerClassName="aspect-square w-full shrink-0 md:w-40"
                />
                <div className="flex min-w-0 flex-1 flex-col justify-between gap-3 p-4">
                    <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
                        <h3 className="min-w-0 flex-1 text-lg font-semibold text-slate-900 md:text-xl">
                            {recipe.name}
                        </h3>
                        {actions && <div className="shrink-0 self-start">{actions}</div>}
                    </div>
                </div>
            </div>
        </div>
    );
}