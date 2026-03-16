import type { Meta, StoryObj } from "@storybook/react-vite";
import { fn } from "storybook/test";
import Button from "./Button";

const meta = {
    title: 'UI/Button',
    component: Button,
    tags: ['autodocs'],
    parameters: { layout: 'centered' },
    args: {
        children: "Save",
        onClick: fn(),
        variant: "primary",
        disabled: false,
        loading: false,
    },
    argTypes: {
        variant: {
            control: { type: 'select' },
            options: ['primary', 'secondary', 'danger', 'primaryOutline', 'secondaryOutline', 'dangerOutline'],
        },
        children: { control: 'text' },
    },
} satisfies Meta<typeof Button>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {};

export const Secondary: Story = {
    args: {
        variant: "secondary",
        children: "Cancel",
    },
};

export const Danger: Story = {
    args: {
        variant: "danger",
        children: "Delete",
    },
};

export const PrimaryOutline: Story = {
    args: {
        variant: "primaryOutline",
        children: "Edit",
    },
};

export const Disabled: Story = {
    args: {
        disabled: true,
        children: "Unavailable",
    },
};

export const Loading: Story = {
    args: {
        loading: true,
        children: "Saving",
    },
};

export const SubmitButton: Story = {
    args: {
        type: "submit",
        children: "Submit",
    },
};