import Stack from "@/components/layout/Stack";
import MarkdownRenderer from "@/components/ui/MarkdownRenderer";
import SectionTitle from "@/components/ui/SectionTitle";
import { RecipeSourceType, type Recipe, type User } from "../types";
import IngredientList from "./IngredientList";
import RecipeDetailsCard from "./RecipeDetailsCard";

type RecipeContentSectionsProps = {
    recipe: Recipe;
    user?: User;
};

function BorderedSection({
    title,
    children,
}: {
    title: string;
    children: React.ReactNode;
}) {
    return (
        <section className="space-y-4 bg-white border border-black rounded-lg p-6 shadow-sm">
            <SectionTitle text={title} />
            {children}
        </section>
    );
}

export default function RecipeContentSections({
    recipe,
    user,
}: RecipeContentSectionsProps) {
    return (
        <Stack space="sm">
            <RecipeDetailsCard
                recipeTitle={recipe.name}
                imageUrl={recipe.imageUrl}
                description={recipe.description}
                prepTimeMinutes={recipe.prepMins}
                cookTimeMinutes={recipe.cookMins}
                portions={recipe.portions}
                versionNumber={recipe.version}
                isDraft={recipe.isDraft}
                sharedBy={user ? { username: user.username, id: user.id } : undefined}
                dietLevel={recipe.animalProductLevel}
                containsGluten={recipe.containsGluten}
            />

            <BorderedSection title="Ingredients">
                <IngredientList ingredients={recipe.ingredients} />
            </BorderedSection>

            {recipe.source.type === RecipeSourceType.Website && recipe.source.url && (
                <BorderedSection title="Website reference">
                    <a
                        href={recipe.source.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-blue-600 hover:underline"
                    >
                        {recipe.source.url}
                    </a>
                </BorderedSection>
            )}

            {recipe.source.type === RecipeSourceType.Cookbook &&
                recipe.source.bookTitle &&
                recipe.source.bookPage && (
                    <BorderedSection title="Book reference">
                        <div>
                            {recipe.source.bookTitle}, page {recipe.source.bookPage}
                        </div>
                    </BorderedSection>
                )}

            {recipe.source.type === RecipeSourceType.Original &&
                recipe.source.instructions && (
                    <BorderedSection title="Instructions">
                        <MarkdownRenderer content={recipe.source.instructions} />
                    </BorderedSection>
                )}
        </Stack>
    );
}
