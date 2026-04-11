import Tag from "@/components/ui/Tag";

export default function DietLevelTag({ level }: { level: number }) {
    let label: string | null;
    switch (level) {
        case 0:
            label = "Vegan";
            break;
        case 1:
            label = "Vegetarian";
            break;
        default:
            label = null;
    }

    return (
        label && <Tag>
            {label}
        </Tag>
    );
}