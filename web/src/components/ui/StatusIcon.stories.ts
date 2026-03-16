import type { Meta, StoryObj } from "@storybook/react-vite";
import StatusIcon from "./StatusIcon";

const meta = {
    title: "UI/StatusIcon",
    component: StatusIcon,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        status: "pending",
    },
    argTypes: {
        status: {
            control: { type: "select" },
            options: ["error", "pending", "completed"],
        },
    },
} satisfies Meta<typeof StatusIcon>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Pending: Story = {};

export const Error: Story = {
    args: {
        status: "error",
    },
};

export const Completed: Story = {
    args: {
        status: "completed",
    },
};