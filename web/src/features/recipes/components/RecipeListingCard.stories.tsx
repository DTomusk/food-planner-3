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
            description: "A classic Italian pasta dish with rich meat sauce.",
            createdAt: "2024-01-01T12:00:00Z",
            author: {
                id: "user1",
                username: "chef123",
            },
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
