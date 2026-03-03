import { Button, Form } from "@/components";
import FormField from "@/components/form/FormField";
import Container from "@/components/layout/Container";
import Input from "@/components/ui/Input";
import { commonStrings } from "@/lib/strings";
import { useForm } from "react-hook-form";

type SignUpFormProps = {
    onSubmit: (values: SignUpFormValues) => void;
    isSubmitting?: boolean;
}

type SignUpFormValues = {
    email: string;
    password: string;
    username: string;
}

export default function SignUpForm({
    onSubmit,
    isSubmitting = false,
}: SignUpFormProps) {
    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<SignUpFormValues>();
    return (
        <Container size="sm">
        <Form onSubmit={handleSubmit(onSubmit)}>
            <FormField htmlFor="email" label="Email" error={errors.email?.message}>
                <Input type="email" placeholder="Email" {...register("email", { required: "Email is required" })} />
            </FormField>
            <FormField htmlFor="username" label="Username" error={errors.username?.message}>
                <Input type="text" placeholder="Username" {...register("username", { required: "Username is required", maxLength: { value: 50, message: "Username cannot exceed 50 characters" } })} />
            </FormField>    
            <FormField htmlFor="password" label="Password" error={errors.password?.message}>
                <Input type="password" placeholder="Password" {...register("password", { required: "Password is required" })} />
            </FormField>      
            <Button disabled={isSubmitting} type="submit" loading={isSubmitting}>
                {commonStrings.auth.signUp}
            </Button>
        </Form>
        </Container>
    );
}