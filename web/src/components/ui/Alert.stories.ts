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