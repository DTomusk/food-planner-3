import type { Meta, StoryObj } from "@storybook/react-vite";
import SectionHelpText from "./SectionHelpText";

const meta = {
    title: "Form/SectionHelpText",
    component: SectionHelpText,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        children: "You can add optional notes for this section.",
    },
    argTypes: {
        children: { control: "text" },
    },
} satisfies Meta<typeof SectionHelpText>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const LongerCopy: Story = {
    args: {
        children: "This section supports markdown-like structure and can include additional guidance for users.",
    },
};