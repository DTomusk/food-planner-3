import { Button, Form } from "@/components";
import FormField from "@/components/form/FormField";
import Container from "@/components/layout/Container";
import Input from "@/components/ui/Input";
import { commonStrings } from "@/lib/strings";
import { useForm } from "react-hook-form";

type SignInFormProps = {
    onSubmit: (values: SignInFormValues) => void;
    isSubmitting?: boolean;
}

type SignInFormValues = {
    email: string;
    password: string;
}

export default function SignInForm({
    onSubmit,
    isSubmitting = false,
}: SignInFormProps) {
    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<SignInFormValues>();
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
                {commonStrings.auth.signIn}
            </Button>
        </Form>
        </Container>
    );
}