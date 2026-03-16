import type { Meta, StoryObj } from "@storybook/react-vite";
import ConfirmModal from "./ConfirmModal";
import { fn } from "storybook/test";

const meta = {
    title: 'UI/ConfirmModal',
    component: ConfirmModal,
    tags: ['autodocs'],
    parameters: { layout: 'centered' },
    args: {
        isOpen: true,
        title: "Are you sure?",
        description: "This action cannot be undone.",
        confirmText: "Delete",
        cancelText: "Cancel",
        loading: false,
        variant: "danger",
        onConfirm: fn(),
        onCancel: fn(),
    },
    argTypes: {
        variant: {
            control: { type: 'select' },
            options: ['primary', 'danger'],
        },
        title: { control: 'text' },
        description: { control: 'text' },
    },
} satisfies Meta<typeof ConfirmModal>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Loading: Story = {
    args: {
        loading: true,
    },
};