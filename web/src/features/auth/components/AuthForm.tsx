// Note: separate this if sign in and up diverge
import { Button, Form } from "@/components";
import FormField from "@/components/form/FormField";
import Container from "@/components/layout/Container";
import Input from "@/components/ui/Input";
import { commonStrings } from "@/lib/strings";
import { useForm } from "react-hook-form";

type AuthFormProps = {
    onSubmit: (values: { email: string; password: string }) => void;
    isSubmitting?: boolean;
    formType: "signin" | "signup";
}

export default function AuthForm({
    onSubmit,
    isSubmitting = false,
    formType,
}: AuthFormProps) {
    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<{ email: string; password: string }>();
    return (
        <Container size="sm">
        <Form onSubmit={handleSubmit(onSubmit)}>
            <FormField htmlFor="email" label="Email" error={errors.email?.message}>
                <Input type="email" placeholder="Email" {...register("email", { required: "Email is required" })} />
            </FormField>
            <FormField htmlFor="password" label="Password" error={errors.password?.message}>
                <Input type="password" placeholder="Password" {...register("password", { required: "Password is required" })} />
            </FormField>      
            <Button disabled={isSubmitting} type="submit" loading={isSubmitting}>
                {formType === "signin" ? commonStrings.auth.signIn : commonStrings.auth.signUp}
            </Button>
        </Form>
        </Container>
    );
}