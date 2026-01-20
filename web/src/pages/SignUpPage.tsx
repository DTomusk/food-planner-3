import { Alert, PageTitle } from "@/components";
import Link from "@/components/Link";
import AuthForm from "@/features/auth/components/AuthForm";
import { useAuth } from "@/features/auth/hooks/useAuth";
import { useSignUp } from "@/features/auth/hooks/useSignUp";
import { Page } from "@/layout";
import { commonStrings } from "@/lib";
import { useNavigate } from "react-router-dom";

export default function SignUpPage() {
    const { mutate, isPending, error } = useSignUp();
    const navigate = useNavigate();
    const { isAuthenticated } = useAuth();

    if (isAuthenticated) {
        navigate("/");
    }

    const handleSubmit = (values: { email: string; password: string }) => {
        mutate(
            { input: { email: values.email, password: values.password } },
            {
                onSuccess: () => {
                    navigate("/");
                }
            }
        );
    }

    return (
        <Page>
            <PageTitle text={commonStrings.auth.signUp} />
            {error && <Alert message={(error as Error).message} />}
            <p>Already have an account? <Link onClick={() => navigate("/auth/signin")} text={commonStrings.auth.signIn}/></p>
            <AuthForm
                formType="signup"
                onSubmit={handleSubmit}
                isSubmitting={isPending}
            />
        </Page>
    )
}