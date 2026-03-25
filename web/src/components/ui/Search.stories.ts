import type { Meta, StoryObj } from "@storybook/react-vite";
import SearchBar from "./SearchBar";

const meta = {
    title: "UI/SearchBar",
    component: SearchBar,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        value: "",
        onChange: (value: string) => console.log("Search value:", value),
        placeholder: "Search...",
        onSubmit: () => console.log("Search submitted"),
        loading: false,
        disabled: false,
    },
    argTypes: {
        value: { control: "text" },
        onChange: { action: "changed" },
        placeholder: { control: "text" },
        onSubmit: { action: "submitted" },
        loading: { control: "boolean" },
        disabled: { control: "boolean" },
    },
} satisfies Meta<typeof SearchBar>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};