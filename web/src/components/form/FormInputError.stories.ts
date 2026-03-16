import type { Meta, StoryObj } from "@storybook/react-vite";
import FormInputError from "./FormInputError";

const meta = {
    title: "Form/FormInputError",
    component: FormInputError,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        message: "This field is required.",
    },
    argTypes: {
        message: { control: "text" },
    },
} satisfies Meta<typeof FormInputError>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Empty: Story = {
    args: {
        message: undefined,
    },
};