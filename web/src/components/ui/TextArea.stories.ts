import type { Meta, StoryObj } from "@storybook/react-vite";
import { fn } from "storybook/test";
import TextArea from "./TextArea";

const meta = {
    title: "UI/TextArea",
    component: TextArea,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        placeholder: "Add notes...",
        rows: 4,
        onChange: fn(),
        disabled: false,
    },
} satisfies Meta<typeof TextArea>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Disabled: Story = {
    args: {
        disabled: true,
        placeholder: "This field is disabled",
    },
};

export const WithValue: Story = {
    args: {
        defaultValue: "Line one\nLine two",
    },
};