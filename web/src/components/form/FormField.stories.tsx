import type { Meta, StoryObj } from "@storybook/react-vite";
import Input from "../ui/Input";
import FormField from "./FormField";

const meta = {
    title: "Form/FormField",
    component: FormField,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        label: "Recipe Name",
        htmlFor: "recipe-name",
        error: undefined,
        children: null,
    },
    argTypes: {
        label: { control: "text" },
        htmlFor: { control: "text" },
        error: { control: "text" },
    },
    render: (args) => (
        <div className="w-80">
            <FormField {...args}>
                <Input id={args.htmlFor} placeholder="Enter a recipe name" />
            </FormField>
        </div>
    ),
} satisfies Meta<typeof FormField>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithError: Story = {
    args: {
        error: "Recipe name is required.",
    },
};