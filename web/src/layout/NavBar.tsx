import Link from "@/components/Link";
import { useAuth } from "@/features/auth/hooks/useAuth";
import { commonStrings } from "@/lib";
import { useNavigate } from "react-router-dom";

export default function NavBar() {
    const navigation = useNavigate();
    const { isAuthenticated, signOut } = useAuth();
    return (
        <nav className="bg-white shadow p-4 mb-6">
            <div className="container mx-auto flex justify-between items-center">
                <h1 className="text-xl font-bold cursor-pointer" onClick={() => navigation("/")}>FoodSmash</h1>
                <div className="space-x-4">
                    <Link onClick={() => navigation("/recipe")} text="Recipes" />
                    {!isAuthenticated && <Link onClick={() => navigation("/auth/signin")} text={commonStrings.auth.signIn} />}
                    {isAuthenticated && <Link onClick={() => signOut()} text={commonStrings.auth.signOut} />}
                </div>
            </div>
        </nav>
    );
}