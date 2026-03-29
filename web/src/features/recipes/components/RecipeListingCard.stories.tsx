import type { Meta, StoryObj } from "@storybook/react-vite";
import RecipeListingCard from "./RecipeListingCard";

const meta = {
    title: "Recipes/RecipeListingCard",
    component: RecipeListingCard,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        recipe: {
            id: "1",
            name: "Spaghetti Bolognese",
            imageUrl: null,
        },
    },
    argTypes: {
        recipe: { control: "object" },
        onClick: { action: "clicked" },
    },
} satisfies Meta<typeof RecipeListingCard>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
