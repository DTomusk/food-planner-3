import { Alert, PageTitle } from "@/components";
import Link from "@/components/Link";
import AuthForm from "@/features/auth/components/AuthForm";
import { useAuth } from "@/features/auth/hooks/useAuth";
import { useSignIn } from "@/features/auth/hooks/useSignIn";
import Page from "@/layout/PageWrapper";
import { commonStrings } from "@/lib";
import { useNavigate } from "react-router-dom";

export default function SignInPage() {
    const { mutate, isPending, error } = useSignIn();
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
            <PageTitle text={commonStrings.auth.signIn} />
            {error && <Alert message={(error as Error).message} />}
            <p>Don't have an account yet? <Link onClick={() => navigate("/auth/signup")} text={commonStrings.auth.signUp}/></p>
            <AuthForm
                formType="signin"
                onSubmit={handleSubmit}
                isSubmitting={isPending}
            />
        </Page>
    );
}