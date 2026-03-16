import type { Meta, StoryObj } from "@storybook/react-vite";
import Text from "./Text";

const meta = {
    title: "UI/Text",
    component: Text,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        children: "This is body text.",
        variant: "body",
        as: "p",
        align: "left",
    },
    argTypes: {
        variant: {
            control: { type: "select" },
            options: ["body", "muted", "caption", "error"],
        },
        as: {
            control: { type: "select" },
            options: ["p", "span", "strong"],
        },
        align: {
            control: { type: "select" },
            options: ["left", "center", "right"],
        },
        children: { control: "text" },
    },
} satisfies Meta<typeof Text>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Body: Story = {};

export const Muted: Story = {
    args: {
        variant: "muted",
        children: "Muted supporting copy.",
    },
};

export const Caption: Story = {
    args: {
        variant: "caption",
        children: "Caption text",
    },
};

export const Error: Story = {
    args: {
        variant: "error",
        children: "This value is invalid.",
    },
};

export const Centered: Story = {
    args: {
        align: "center",
        children: "Centered text",
    },
};