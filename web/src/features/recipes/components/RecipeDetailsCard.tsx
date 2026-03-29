import { Inline, Stack } from "@/components";
import SharedBy from "@/components/SharedBy";
import ImageDisplay from "@/components/ui/ImageDisplay";

type RecipeDetailsCardProps = {
    recipeTitle: string;
    imageUrl: string;
    description: string;
    prepTimeMinutes: number;
    cookTimeMinutes: number;
    portions: number;
    sharedBy?: {
        username: string;
        id: string;
    }
};

export default function RecipeDetailsCard({ 
    recipeTitle, 
    imageUrl, 
    description, 
    prepTimeMinutes, 
    cookTimeMinutes, 
    portions,
    sharedBy
}: RecipeDetailsCardProps) {
    return (
        <div className="flex flex-col border border-black bg-white rounded shadow">
        <div className="flex flex-col md:flex-row border-b border-black rounded-t">
            <ImageDisplay
                imageUrl={imageUrl}
                altText={recipeTitle}
                containerClassName="aspect-square w-full shrink-0 md:w-64 overflow-hidden"
            />
            <Stack className="py-3 px-4">
                <h2 className="text-3xl font-bold">{recipeTitle}</h2>
                {sharedBy && <SharedBy user={sharedBy} />}
                <p className="text-gray-700">{description}</p>
                
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