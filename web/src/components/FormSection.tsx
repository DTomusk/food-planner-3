import FormSectionTitle from "./FormSectionTitle";

interface FormSectionProps {
    title: string;
    children: React.ReactNode;
}

export default function FormSection({ title, children }: FormSectionProps) {
    return (
        <section className="space-y-4">
            <FormSectionTitle title={title} />
            {children}
        </section>
    );
}