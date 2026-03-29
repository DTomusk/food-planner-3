import { Home } from "lucide-react";
import NavItem from "./NavItem";
import type { Meta, StoryObj } from "@storybook/react-vite";

const meta = {
    title: "Nav/NavItem",
    component: NavItem,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        label: "Home",
        icon: <Home />,
    },
    argTypes: {
        label: { control: "text" },
        icon: { control: "boolean" },
    },
} satisfies Meta<typeof NavItem>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const NoIcon: Story = {
    args: {
        icon: undefined,
    },
};