import FormSectionTitle from "./FormSectionTitle";

interface FormSectionProps {
    title?: string;
    children: React.ReactNode;
}

export default function FormSection({ title, children }: FormSectionProps) {
    return (
        <section className="space-y-4">
            {title && <FormSectionTitle title={title} />}
            {children}
        </section>
    );
}