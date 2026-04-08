import type { Meta, StoryObj } from "@storybook/react-vite";
import ImageDisplay from "./ImageDisplay";

const meta = {
    title: "UI/ImageDisplay",
    component: ImageDisplay,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        imageUrl: "https://placehold.co/600x400?text=Recipe",
        altText: "Placeholder Image",
    },
    argTypes: {
        imageUrl: { control: "text" },
        altText: { control: "text" },
    },
} satisfies Meta<typeof ImageDisplay>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const BrokenImage: Story = {
    args: {
        imageUrl: "https://example.com/nonexistent.jpg",
    },
};