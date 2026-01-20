import { Alert, PageTitle } from "@/components";
import AuthForm from "@/features/auth/components/AuthForm";
import { useSignUp } from "@/features/auth/hooks/useSignUp";
import { Page } from "@/layout";

export default function SignUpPage() {
    const { mutate, isPending, error } = useSignUp();

    const handleSubmit = (values: { email: string; password: string }) => {
        mutate(
            { input: { email: values.email, password: values.password } },
            {
                // TODO: store jwt and redirect
                onSuccess: () => {
                    alert("Sign up successful!");
                }
            }
        );
    }

    return (
        <Page>
            <PageTitle text="Sign Up" />
            <AuthForm
                formType="signup"
                onSubmit={handleSubmit}
                isSubmitting={isPending}
            />
            {error && <Alert message={(error as Error).message} />}
        </Page>
    )
}