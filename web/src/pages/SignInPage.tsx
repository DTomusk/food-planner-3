import { Alert, PageTitle } from "@/components";
import Link from "@/components/ui/Link";
import SignInForm from "@/features/auth/components/SignInForm";
import { useAuth } from "@/features/auth/hooks/useAuth";
import { useSignIn } from "@/features/auth/hooks/useSignIn";
import Page from "@/layout/PageWrapper";
import { commonStrings } from "@/lib";
import { extractErrorMessage } from "@/lib/errors";
import { Navigate, useLocation, useNavigate } from "react-router-dom";

export default function SignInPage() {
    const { mutate, isPending, error } = useSignIn();
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
            <PageTitle text={commonStrings.auth.signIn} />
            {error && <Alert message={extractErrorMessage(error)} closable />}
            <p className="text-center">Don't have an account yet? <Link onClick={() => navigate("/auth/signup", { state: { from: location } })} text={commonStrings.auth.signUp}/></p>
            <SignInForm
                onSubmit={handleSubmit}
                isSubmitting={isPending}
            />
        </Page>
    );
}