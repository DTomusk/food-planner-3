import type { Meta, StoryObj } from "@storybook/react-vite";
import { Ellipsis } from "lucide-react";
import { fn } from "storybook/test";
import IconButton from "./IconButton";

const meta = {
    title: "UI/IconButton",
    component: IconButton,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        "aria-label": "More actions",
        onClick: fn(),
        variant: "primary",
        shape: "square",
        disabled: false,
        children: <Ellipsis size={16} />,
    },
    argTypes: {
        variant: {
            control: { type: "select" },
            options: ["primary", "secondary", "danger", "primary-outline", "secondary-outline"],
        },
        shape: {
            control: { type: "select" },
            options: ["circle", "square"],
        },
        children: { control: false },
    },
} satisfies Meta<typeof IconButton>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {};

export const Secondary: Story = {
    args: {
        variant: "secondary",
    },
};

export const Danger: Story = {
    args: {
        variant: "danger",
    },
};

export const Circle: Story = {
    args: {
        shape: "circle",
    },
};

export const PrimaryOutline: Story = {
    args: {
        variant: "primary-outline",
    },
};

export const Disabled: Story = {
    args: {
        disabled: true,
    },
};