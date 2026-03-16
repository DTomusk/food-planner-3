import type { Meta, StoryObj } from "@storybook/react-vite";
import FormSectionTitle from "./FormSectionTitle";

const meta = {
    title: "Form/FormSectionTitle",
    component: FormSectionTitle,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        title: "Ingredients",
    },
    argTypes: {
        title: { control: "text" },
    },
} satisfies Meta<typeof FormSectionTitle>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const LongTitle: Story = {
    args: {
        title: "Preparation and Timing",
    },
};