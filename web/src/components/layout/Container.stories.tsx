import type { Meta, StoryObj } from "@storybook/react-vite";
import Container from "./Container";

const demoContent = (
    <div className="rounded-md border border-dashed border-gray-400 bg-white p-4 text-sm text-gray-700">
        Container content
    </div>
);

const meta = {
    title: "Layout/Container",
    component: Container,
    tags: ["autodocs"],
    parameters: { layout: "fullscreen" },
    args: {
        size: "md",
        children: demoContent,
    },
    argTypes: {
        size: {
            control: { type: "select" },
            options: ["xs", "sm", "md", "lg", "xl", "full"],
        },
        children: { control: false },
    },
    render: (args) => (
        <div className="w-full bg-gray-50 py-8">
            <Container {...args} />
        </div>
    ),
} satisfies Meta<typeof Container>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Medium: Story = {};

export const ExtraSmall: Story = {
    args: {
        size: "xs",
    },
};

export const Large: Story = {
    args: {
        size: "lg",
    },
};

export const ExtraLarge: Story = {
    args: {
        size: "xl",
    },
};

export const FullWidth: Story = {
    args: {
        size: "full",
    },
};