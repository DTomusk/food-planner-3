import type { Meta, StoryObj } from "@storybook/react-vite";
import { fn } from "storybook/test";
import Link from "./Link";

const meta = {
    title: "UI/Link",
    component: Link,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        children: "View recipe details",
        color: "black",
        onClick: fn(),
    },
    argTypes: {
        color: {
            control: { type: "select" },
            options: ["black", "primary"],
        },
        children: { control: "text" },
    },
} satisfies Meta<typeof Link>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Black: Story = {};

export const Primary: Story = {
    args: {
        color: "primary",
        children: "Go back",
    },
};