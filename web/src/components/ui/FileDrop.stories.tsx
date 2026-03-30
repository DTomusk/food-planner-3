import type { Meta, StoryObj } from "@storybook/react-vite";
import FileDrop from "./FileDrop";

const meta = {
    title: "UI/FileDrop",
    component: FileDrop,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        value: null,
        onChange: (file: File | null) => console.log("File changed:", file),
        accept: { "image/*": [] },
        maxSize: 5 * 1024 * 1024, // 5 MB
    },
    argTypes: {
        value: { control: "object" },
        onChange: { action: "changed" },
        accept: { control: "object" },
        maxSize: { control: "number" },
    },
} satisfies Meta<typeof FileDrop>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};