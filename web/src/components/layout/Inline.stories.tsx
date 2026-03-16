import type { Meta, StoryObj } from "@storybook/react-vite";
import Inline from "./Inline";

const inlineContent = (
    <>
        <div className="min-w-20 rounded bg-blue-100 px-3 py-2 text-center text-xs font-medium text-blue-900">One</div>
        <div className="min-w-20 rounded bg-green-100 px-3 py-2 text-center text-xs font-medium text-green-900">Two</div>
        <div className="min-w-20 rounded bg-amber-100 px-3 py-2 text-center text-xs font-medium text-amber-900">Three</div>
    </>
);

const wrappedContent = (
    <>
        <div className="min-w-24 rounded bg-blue-100 px-3 py-2 text-center text-xs font-medium text-blue-900">Alpha</div>
        <div className="min-w-24 rounded bg-green-100 px-3 py-2 text-center text-xs font-medium text-green-900">Beta</div>
        <div className="min-w-24 rounded bg-amber-100 px-3 py-2 text-center text-xs font-medium text-amber-900">Gamma</div>
        <div className="min-w-24 rounded bg-red-100 px-3 py-2 text-center text-xs font-medium text-red-900">Delta</div>
        <div className="min-w-24 rounded bg-purple-100 px-3 py-2 text-center text-xs font-medium text-purple-900">Epsilon</div>
    </>
);

const meta = {
    title: "Layout/Inline",
    component: Inline,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        justify: "center",
        align: "none",
        wrap: false,
        gap: "sm",
        children: inlineContent,
    },
    argTypes: {
        justify: {
            control: { type: "select" },
            options: ["start", "center", "end", "between", "around", "evenly"],
        },
        align: {
            control: { type: "select" },
            options: ["start", "center", "end", "none"],
        },
        gap: {
            control: { type: "select" },
            options: ["none", "sm", "md", "lg"],
        },
        children: { control: false },
    },
    render: (args) => (
        <div className="w-[560px] rounded border border-dashed border-gray-300 bg-white p-4">
            <Inline {...args} />
        </div>
    ),
} satisfies Meta<typeof Inline>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const JustifyBetween: Story = {
    args: {
        justify: "between",
    },
};

export const JustifyEvenly: Story = {
    args: {
        justify: "evenly",
    },
};

export const AlignCenter: Story = {
    args: {
        align: "center",
    },
};

export const Wrapped: Story = {
    args: {
        justify: "start",
        align: "start",
        wrap: true,
        children: wrappedContent,
    },
};

export const LargeGap: Story = {
    args: {
        justify: "start",
        gap: "lg",
    },
};

export const NoGap: Story = {
    args: {
        justify: "start",
        gap: "none",
    },
};