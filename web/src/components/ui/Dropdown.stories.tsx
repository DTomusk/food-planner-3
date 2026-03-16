import { fn } from "storybook/test";
import Dropdown from "./Dropdown";
import type { Meta, StoryObj } from "@storybook/react-vite";
import Button from "./Button";
import { Clock, PersonStanding } from "lucide-react";

const meta = {
    title: "UI/Dropdown",
    component: Dropdown,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        button: <Button>Open Dropdown</Button>,
        sections: [
            {
                title: "Section 1",
                items: [
                    { label: "Item 1", onClick: fn() },
                    { label: "Item 2", onClick: fn() },
                ],
            },
            {
                title: "Section 2",
                items: [
                    { label: "Danger Item", onClick: fn(), danger: true },
                    { label: "Disabled Item", onClick: fn(), disabled: true },
                ],
            },
        ],
    },
    argTypes: {
        button: {
            control: false,
        },
    },
} satisfies Meta<typeof Dropdown>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithIcons: Story = {
    args: {
        sections: [
            {
                title: "Actions",
                items: [
                    { label: "Edit", onClick: fn(), icon: <Clock /> },
                    { label: "Duplicate", onClick: fn(), icon: <PersonStanding /> },
                ],
            },
        ],
    },
};

export const LotsOfItems: Story = {
    args: {
        sections: [
            {
                items: Array.from({ length: 20 }, (_, i) => ({
                    label: `Item ${i + 1}`,
                    onClick: fn(),
                })),
            },
        ],
    },
};