import type { Meta, StoryObj } from "@storybook/react-vite";
import ResizableSidebar from "./ResizableSidebar";

const defaultChildren = (
    <>
        <div className="text-lg font-semibold tracking-tight">FoodSmash</div>
        <div className="mt-6 space-y-3 text-sm">
            <div>Home</div>
            <div>My recipes</div>
            <div>Sign out</div>
        </div>
    </>
);

const meta = {
    title: "Layout/ResizableSidebar",
    component: ResizableSidebar,
    tags: ["autodocs"],
    parameters: { layout: "fullscreen" },
    args: {
        children: defaultChildren,
        minWidth: 192,
        maxWidth: 384,
        defaultWidth: 256,
        className: "h-screen border-r border-black bg-white px-4 py-6",
    },
    argTypes: {
        minWidth: { control: { type: "number", min: 120, max: 420, step: 8 } },
        maxWidth: { control: { type: "number", min: 200, max: 560, step: 8 } },
        defaultWidth: { control: { type: "number", min: 120, max: 560, step: 8 } },
        className: { control: "text" },
    },
    render: (args) => (
        <div className="min-h-screen min-w-[900px] bg-linear-to-br from-background to-white">
            <div className="flex min-h-screen w-full">
                <ResizableSidebar
                    key={`${args.minWidth ?? 192}-${args.maxWidth ?? 384}-${args.defaultWidth ?? 256}`}
                    minWidth={args.minWidth}
                    maxWidth={args.maxWidth}
                    defaultWidth={args.defaultWidth}
                    className={args.className}
                >
                    {args.children}
                </ResizableSidebar>
                <main className="flex-1 p-8">
                    <h2 className="text-xl font-semibold">Main Content</h2>
                    <p className="mt-2 text-sm text-muted-foreground">
                        Drag the sidebar handle at its right edge to resize.
                    </p>
                </main>
            </div>
        </div>
    ),
} satisfies Meta<typeof ResizableSidebar>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const NarrowRange: Story = {
    args: {
        minWidth: 208,
        maxWidth: 280,
        defaultWidth: 240,
    },
};

export const WideRange: Story = {
    args: {
        minWidth: 176,
        maxWidth: 448,
        defaultWidth: 320,
    },
};
