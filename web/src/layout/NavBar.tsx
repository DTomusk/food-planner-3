import Stack from "@/components/layout/Stack";
import NavItem from "@/components/ui/NavItem";
import { useAuth } from "@/features/auth/hooks/useAuth";
import { useSignOut } from "@/features/auth/hooks/useSignOut";
import { commonStrings } from "@/lib";
import { BookOpen, Home, LogIn, LogOut } from "lucide-react";
import { useLocation, useNavigate } from "react-router-dom";

export default function NavBar() {
    const navigate = useNavigate();
    const { mutate } = useSignOut();

    const location = useLocation();
    
    const handleSignInClick = () => {
        navigate("/auth/signin", {
            state: { from: location },
        });
    };

    const { isAuthenticated } = useAuth();
    return (
        <nav className="w-full shrink-0 border-b border-black bg-white px-4 py-4 sm:sticky sm:top-0 sm:h-screen sm:w-64 sm:border-b-0 sm:border-r sm:py-6">
            <Stack space="lg">
                <h1 className="cursor-pointer text-lg font-semibold tracking-tight" onClick={() => navigate("/")}>FoodSmash</h1>
                <Stack space="md">
                    <NavItem onClick={() => navigate("/")} icon={<Home />} label="Home" />
                    {isAuthenticated && <NavItem onClick={() => navigate("/me/recipes")} icon={<BookOpen />} label="My recipes" />}
                    {!isAuthenticated && <NavItem onClick={handleSignInClick} icon={<LogIn />} label={commonStrings.auth.signIn} />}
                    {isAuthenticated && <NavItem onClick={() => mutate()} icon={<LogOut />} label={commonStrings.auth.signOut} />}
                </Stack>
            </Stack>
        </nav>
    );
}