import type { Meta, StoryObj } from "@storybook/react-vite";
import SectionTitle from "./SectionTitle";

const meta = {
    title: "UI/SectionTitle",
    component: SectionTitle,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        text: "Ingredients",
    },
    argTypes: {
        text: { control: "text" },
    },
} satisfies Meta<typeof SectionTitle>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const AlternateText: Story = {
    args: {
        text: "Preparation",
    },
};