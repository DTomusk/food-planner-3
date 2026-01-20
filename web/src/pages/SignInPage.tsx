import { Alert, PageTitle } from "@/components";
import AuthForm from "@/features/auth/components/AuthForm";
import { useSignIn } from "@/features/auth/hooks/useSignIn";
import Page from "@/layout/PageWrapper";

export default function SignInPage() {
    const { mutate, isPending, error } = useSignIn();

    const handleSubmit = (values: { email: string; password: string }) => {
        mutate(
            { input: { email: values.email, password: values.password } },
            {
                onSuccess: () => {
                    alert("Sign in successful!");
                }
            }
        );
    }
            
    return (
        <Page>
            <PageTitle text="Sign In" />
            {error && <Alert message={(error as Error).message} />}
            <AuthForm
                formType="signin"
                onSubmit={handleSubmit}
                isSubmitting={isPending}
            />
        </Page>
    );
}