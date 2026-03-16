import type { Meta, StoryObj } from "@storybook/react-vite";
import FormSection from "./FormSection";

const meta = {
    title: "Form/FormSection",
    component: FormSection,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        title: "Recipe Details",
        collapsible: false,
        defaultCollapsed: false,
        showStatus: false,
        status: { status: "pending" },
    },
    argTypes: {
        title: { control: "text" },
        status: { control: false },
    },
    render: (args) => (
        <div className="w-[520px]">
            <FormSection {...args}>
                <p className="text-sm text-gray-700">Add fields and helper text for this section.</p>
            </FormSection>
        </div>
    ),
} satisfies Meta<typeof FormSection>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Collapsible: Story = {
    args: {
        collapsible: true,
    },
};

export const DefaultCollapsed: Story = {
    args: {
        collapsible: true,
        defaultCollapsed: true,
    },
};

export const WithPendingStatus: Story = {
    args: {
        showStatus: true,
        status: { status: "pending" },
    },
};

export const WithCompletedStatus: Story = {
    args: {
        showStatus: true,
        status: { status: "completed" },
    },
};

export const WithErrorStatus: Story = {
    args: {
        showStatus: true,
        status: { status: "error" },
    },
};