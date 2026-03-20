import type { Meta, StoryObj } from "@storybook/react-vite";
import SearchDropdown from "./SearchDropdown";
import { useState } from "react";

const meta = {
    title: "UI/SearchDropdown",
    component: SearchDropdown,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        items: [
            { label: "Apple", value: "apple" },
            { label: "Banana", value: "banana" },
            { label: "Cherry", value: "cherry" },
        ],
        onSelect: (item: { label: string; value: string } | null) => console.log("Selected:", item),
        selectedItem: null,
    },
    argTypes: {
        items: { control: "object" },
        onSelect: { action: "selected" },
        selectedItem: { control: "object" },
    },
} satisfies Meta<typeof SearchDropdown>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
    render: (args) => {
        const [selectedItem, setSelectedItem] = useState<{ label: string; value: string } | null>(null);
        return (
            <SearchDropdown 
                {...args} 
                selectedItem={selectedItem} 
                onSelect={(item) => {
                    setSelectedItem(item);
                    args.onSelect(item);
                }} 
            />
        );
    }
};

export const EmptyItems: Story = {
    args: {
        items: [],
    },
};