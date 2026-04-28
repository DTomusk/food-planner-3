import { Inline, Stack, Text, } from "@/components";
import type { RecipeSummary } from "../types";
import ImageDisplay from "@/components/ui/ImageDisplay";
import Heading from "@/components/ui/Heading";
import SharedBy from "@/components/SharedBy";
import DietLevelTag from "./DietLevelTag";
import Tag from "@/components/ui/Tag";

interface RecipeListingCardProps {
    recipe: RecipeSummary;
    onClick?: () => void;
    actions?: React.ReactNode;
    dietLevel: number;
}

export default function RecipeListingCard({ recipe, onClick, actions, dietLevel }: RecipeListingCardProps) {
    return (
        <div
            className="cursor-pointer overflow-hidden rounded border bg-white shadow transition-transform hover:scale-[1.02]"
            onClick={onClick}
        >
            <div className="flex flex-col md:flex-row">
                <ImageDisplay
                    imageUrl={recipe.imageUrl}
                    altText={recipe.name}
                    containerClassName="aspect-square w-full shrink-0 md:w-48 lg:w-64"
                />
                <div className="flex min-w-0 flex-1 flex-col justify-between gap-3 p-4">
                    <Stack space="sm">
                    <Tag>Version {recipe.version}</Tag>
                    <Inline justify="between" align="start" gap="md">
                        <Heading text={recipe.name} />
                        {actions && <Inline className="shrink-0 self-start">{actions}</Inline>}
                    </Inline>
                    <SharedBy user={recipe.author} />
                    <div className="hidden md:block">
                        <Text variant="muted">Created on {new Date(recipe.createdAt).toLocaleDateString()}</Text>
                    </div>
                    <Text variant="body">{recipe.description}</Text>
                    <Inline justify="start">
                        <DietLevelTag level={dietLevel} /> 
                        {!recipe.containsGluten && <Tag>Gluten free</Tag>}
                    </Inline>
                    </Stack>
                </div>
            </div>
        </div>
    );
}