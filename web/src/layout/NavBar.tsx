import Link from "@/components/ui/Link";
import { useAuth } from "@/features/auth/hooks/useAuth";
import { commonStrings } from "@/lib";
import { useLocation, useNavigate } from "react-router-dom";

export default function NavBar() {
    const navigate = useNavigate();

    const location = useLocation();
    
    const handleSignInClick = () => {
        navigate("/auth/signin", {
            state: { from: location },
        });
    };

    const { isAuthenticated, signOut } = useAuth();
    return (
        <nav className="bg-white shadow p-4 mb-6">
            <div className="container mx-auto flex justify-between items-center">
                <h1 className="text-xl font-bold cursor-pointer" onClick={() => navigate("/")}>FoodSmash</h1>
                <div className="space-x-4">
                    <Link onClick={() => navigate("/me/recipes")}>My recipes</Link>
                    {!isAuthenticated && <Link onClick={handleSignInClick}>{commonStrings.auth.signIn}</Link>}
                    {isAuthenticated && <Link onClick={() => signOut()}>{commonStrings.auth.signOut}</Link>}
                </div>
            </div>
        </nav>
    );
}