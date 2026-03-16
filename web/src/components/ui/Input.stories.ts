import type { Meta, StoryObj } from "@storybook/react-vite";
import { fn } from "storybook/test";
import Input from "./Input";

const meta = {
    title: "UI/Input",
    component: Input,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        placeholder: "Enter text",
        onChange: fn(),
        disabled: false,
    },
    argTypes: {
        error: { control: "text" },
    },
} satisfies Meta<typeof Input>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithError: Story = {
    args: {
        error: "This field is required.",
    },
};

export const Disabled: Story = {
    args: {
        disabled: true,
        placeholder: "Unavailable input",
    },
};

export const WithValue: Story = {
    args: {
        defaultValue: "Prefilled value",
    },
};