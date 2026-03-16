import type { Meta, StoryObj } from "@storybook/react-vite";
import { fn } from "storybook/test";
import Select from "./Select";

const mealTypeOptions = (
    <>
        <option value="">Select a meal type</option>
        <option value="breakfast">Breakfast</option>
        <option value="lunch">Lunch</option>
        <option value="dinner">Dinner</option>
    </>
);

const meta = {
    title: "UI/Select",
    component: Select,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        children: mealTypeOptions,
        defaultValue: "",
        onChange: fn(),
        disabled: false,
    },
    argTypes: {
        children: { control: false },
        error: { control: "text" },
    },
} satisfies Meta<typeof Select>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithValue: Story = {
    args: {
        defaultValue: "dinner",
    },
};

export const WithError: Story = {
    args: {
        error: "Please select an option.",
    },
};

export const Disabled: Story = {
    args: {
        disabled: true,
        defaultValue: "lunch",
    },
};