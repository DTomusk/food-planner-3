import type { Meta, StoryObj } from "@storybook/react-vite";
import PageTitle from "./PageTitle";

const meta = {
    title: "UI/PageTitle",
    component: PageTitle,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        text: "Recipe Planner",
    },
    argTypes: {
        text: { control: "text" },
    },
} satisfies Meta<typeof PageTitle>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const LongTitle: Story = {
    args: {
        text: "Weekly Family Meal Planning",
    },
};