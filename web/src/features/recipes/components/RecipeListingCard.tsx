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
            <div className="flex flex-col sm:flex-row">
                <ImageDisplay
                    imageUrl={recipe.imageUrl}
                    altText={recipe.name}
                    containerClassName="aspect-square w-full shrink-0 sm:w-40 md:w-48"
                />
                <div className="flex min-w-0 flex-1 flex-col justify-between gap-3 p-4">
                    <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                        <h3 className="min-w-0 flex-1 text-lg font-semibold text-slate-900 sm:text-xl">
                            {recipe.name}
                        </h3>
                        {actions && <div className="shrink-0 self-start">{actions}</div>}
                    </div>
                </div>
            </div>
        </div>
    );
}