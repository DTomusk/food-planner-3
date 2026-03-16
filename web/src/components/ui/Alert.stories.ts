import { fn } from "storybook/test";
import Alert from "./Alert";
import type { Meta, StoryObj } from "@storybook/react-vite";

const meta = {
    title: 'UI/Alert',
    component: Alert,
    tags: ['autodocs'],
    parameters: { layout: 'centered' },
    args: {
        message: "This is an alert message.",
        type: "error",
        closable: true,
        onClose: fn(),
    },
    argTypes: {
        type: {
            control: { type: 'select' },
            options: ['info', 'error', 'success', 'warning'],
        },
        message: { control: 'text' },
        duration: { control: 'number' },
    },
} satisfies Meta<typeof Alert>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Error: Story = {};

export const Success: Story = {
    args: {
        type: "success",
        message: "Operation completed successfully!",
    },
};

export const Warning: Story = {
    args: {
        type: "warning",
        message: "This is a warning. Please be cautious.",
    },
};

export const Info: Story = {
    args: {
        type: "info",
        message: "This is some informational message.",
    },
};

export const AutoClose: Story = {
    args: {
        message: "This alert will auto-close after 3 seconds.",
        duration: 3000,
    },
};

export const NonClosable: Story = {
    args: {
        message: "This alert cannot be closed manually.",
        closable: false,
        type: "warning",
    },
};