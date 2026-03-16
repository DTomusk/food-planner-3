import type { Meta, StoryObj } from "@storybook/react-vite";
import CheckListItem from "./CheckListItem";

const meta = {
    title: 'UI/CheckListItem',
    component: CheckListItem,
    tags: ['autodocs'],
    parameters: { layout: 'centered' },
    args: {
        checked: false,
        onChange: () => {},
        children: "Sample Item",
    },
    argTypes: {
        checked: { control: 'boolean' },
        children: { control: 'text' },
    },
} satisfies Meta<typeof CheckListItem>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Checked: Story = {
    args: {
        checked: true,
    },
};
