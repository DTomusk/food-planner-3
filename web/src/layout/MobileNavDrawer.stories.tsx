import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { BookOpen, Home, LogOut } from "lucide-react";
import { fn } from "storybook/test";
import MobileNavDrawer from "./MobileNavDrawer";
import NavItem from "./NavItem";
import Stack from "@/components/layout/Stack";

function DefaultChildren() {
    return (
        <Stack space="md">
            <NavItem icon={<Home />} label="Home" />
            <NavItem icon={<BookOpen />} label="My recipes" />
            <NavItem icon={<LogOut />} label="Sign out" />
        </Stack>
    );
}

const meta = {
    title: "Nav/MobileNavDrawer",
    component: MobileNavDrawer,
    tags: ["autodocs"],
    parameters: { layout: "fullscreen" },
    args: {
        open: true,
        onClose: fn(),
        children: <DefaultChildren />,
    },
    argTypes: {
        open: { control: "boolean" },
    },
} satisfies Meta<typeof MobileNavDrawer>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Open: Story = {};

export const Interactive: Story = {
    render: () => {
        const [open, setOpen] = useState(false);
        return (
            <div className="p-4">
                <button
                    className="rounded-md bg-primary-600 px-4 py-2 text-sm font-medium text-white hover:bg-primary-700"
                    onClick={() => setOpen(true)}
                >
                    Open drawer
                </button>
                <MobileNavDrawer open={open} onClose={() => setOpen(false)}>
                    <DefaultChildren />
                </MobileNavDrawer>
            </div>
        );
    },
};
