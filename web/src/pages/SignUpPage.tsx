import { Alert, PageTitle } from "@/components";
import Link from "@/components/ui/Link";
import AuthForm from "@/features/auth/components/AuthForm";
import { useAuth } from "@/features/auth/hooks/useAuth";
import { useSignUp } from "@/features/auth/hooks/useSignUp";
import { Page } from "@/layout";
import { commonStrings } from "@/lib";
import { extractErrorMessage } from "@/lib/errors";
import { Navigate, useLocation, useNavigate } from "react-router-dom";

export default function SignUpPage() {
    const { mutate, isPending, error } = useSignUp();
    const navigate = useNavigate();
    const { isAuthenticated } = useAuth();

    const location = useLocation();
    const from = location.state?.from?.pathname || "/";

    if (isAuthenticated) {
        return (<Navigate to={from} />);
    }

    const handleSubmit = (values: { email: string; password: string }) => {
        mutate(
            { input: { email: values.email, password: values.password } },
            {
                onSuccess: () => {
                    navigate(from, { replace: true });
                }
            }
        );
    }

    return (
        <Page>
            <PageTitle text={commonStrings.auth.signUp} />
            {error && <Alert message={extractErrorMessage(error)} closable />}
            <p className="text-center">Already have an account? <Link onClick={() => navigate("/auth/signin", { state: { from: location } })} text={commonStrings.auth.signIn}/></p>
            <AuthForm
                formType="signup"
                onSubmit={handleSubmit}
                isSubmitting={isPending}
            />
        </Page>
    )
}