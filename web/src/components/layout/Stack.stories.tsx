import type { Meta, StoryObj } from "@storybook/react-vite";
import Stack from "./Stack";

const stackContent = (
    <>
        <div className="rounded bg-blue-100 p-3 text-sm font-medium text-blue-900">First block</div>
        <div className="rounded bg-green-100 p-3 text-sm font-medium text-green-900">Second block</div>
        <div className="rounded bg-amber-100 p-3 text-sm font-medium text-amber-900">Third block</div>
    </>
);

const meta = {
    title: "Layout/Stack",
    component: Stack,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        space: "md",
        children: stackContent,
    },
    argTypes: {
        space: {
            control: { type: "select" },
            options: ["sm", "md", "lg", "xl"],
        },
        children: { control: false },
    },
    render: (args) => (
        <div className="w-80 rounded border border-gray-200 bg-white p-4">
            <Stack {...args} />
        </div>
    ),
} satisfies Meta<typeof Stack>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Medium: Story = {};

export const Small: Story = {
    args: {
        space: "sm",
    },
};

export const Large: Story = {
    args: {
        space: "lg",
    },
};

export const ExtraLarge: Story = {
    args: {
        space: "xl",
    },
};