import type { Meta, StoryObj } from "@storybook/react-vite";
import { fn } from "storybook/test";
import Input from "../ui/Input";
import Button from "../ui/Button";
import Form from "./FormWrapper";

const meta = {
    title: "Form/FormWrapper",
    component: Form,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        onSubmit: fn(),
        children: null,
    },
    argTypes: {
        onSubmit: { control: false },
    },
    render: (args) => (
        <div className="w-96">
            <Form
                onSubmit={(event) => {
                    event.preventDefault();
                    args.onSubmit(event);
                }}
            >
                <div className="space-y-3">
                    <Input id="recipe-name" name="recipeName" placeholder="Recipe name" />
                    <Button type="submit">Submit</Button>
                </div>
            </Form>
        </div>
    ),
} satisfies Meta<typeof Form>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};