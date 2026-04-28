import { Inline, Stack } from "@/components";
import SharedBy from "@/components/SharedBy";
import ImageDisplay from "@/components/ui/ImageDisplay";
import Tag from "@/components/ui/Tag";
import DietLevelTag from "./DietLevelTag";

type RecipeDetailsCardProps = {
    recipeTitle: string;
    imageUrl: string | null;
    description: string;
    prepTimeMinutes: number;
    cookTimeMinutes: number;
    portions: number;
    versionNumber: number;
    sharedBy?: {
        username: string;
        id: string;
    },
    dietLevel: number;
    containsGluten: boolean;
    isDraft: boolean;
};

export default function RecipeDetailsCard({ 
    recipeTitle, 
    imageUrl, 
    description, 
    prepTimeMinutes, 
    cookTimeMinutes, 
    portions,
    versionNumber,
    sharedBy,
    dietLevel,
    containsGluten,
    isDraft,
}: RecipeDetailsCardProps) {
    return (
        <div className="flex flex-col border border-black bg-white rounded shadow">
            <div className="flex flex-col md:flex-row border-b border-black rounded-t overflow-hidden">
                <ImageDisplay
                    imageUrl={imageUrl}
                    altText={recipeTitle}
                    containerClassName="aspect-square w-full shrink-0 md:w-64 overflow-hidden"
                />
                <Stack className="py-3 px-4">
                    <Tag>Version {versionNumber} {isDraft ? "(Draft)" : ""}</Tag>
                    <h2 className="text-3xl font-bold">{recipeTitle}</h2>
                    {sharedBy && <SharedBy user={sharedBy} />}
                    <p className="text-gray-700">{description}</p>
                    <Inline justify="start" gap="sm">
                        <DietLevelTag level={dietLevel} /> 
                        {!containsGluten && <Tag>Gluten free</Tag>}
                    </Inline>
                </Stack>
            </div>
            <div className="bg-primary-600 text-white rounded-b px-4 py-3 mt-[-1px]">
                <Inline justify="center" gap="lg">
                    <span>Prep Time: {prepTimeMinutes} mins</span>
                    <span>Cook Time: {cookTimeMinutes} mins</span>
                    <span>Portions: {portions}</span>
                </Inline>
            </div>
        </div>
    );
}