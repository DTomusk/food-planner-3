import { Stack } from "@/components";
import NavItem from "@/layout/NavItem";
import IconButton from "@/components/ui/IconButton";
import MobileNavDrawer from "./MobileNavDrawer";
import { useAuth } from "@/features/auth/hooks/useAuth";
import { useSignOut } from "@/features/auth/hooks/useSignOut";
import { commonStrings } from "@/lib";
import { BookOpen, Home, LogIn, LogOut, Menu } from "lucide-react";
import { useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";

export default function NavBar() {
    const navigate = useNavigate();
    const { mutate } = useSignOut();
    const location = useLocation();
    const { isAuthenticated } = useAuth();
    const [drawerOpen, setDrawerOpen] = useState(false);

    const handleSignInClick = () => {
        navigate("/auth/signin", { state: { from: location } });
    };

    const handleNav = (fn: () => void) => () => {
        fn();
        setDrawerOpen(false);
    };

    const navItems = (
        <Stack space="md">
            <NavItem onClick={handleNav(() => navigate("/"))} icon={<Home />} label="Home" />
            {isAuthenticated && <NavItem onClick={handleNav(() => navigate("/me/recipes"))} icon={<BookOpen />} label="My recipes" />}
            {!isAuthenticated && <NavItem onClick={handleNav(handleSignInClick)} icon={<LogIn />} label={commonStrings.auth.signIn} />}
            {isAuthenticated && <NavItem onClick={handleNav(() => mutate())} icon={<LogOut />} label={commonStrings.auth.signOut} />}
        </Stack>
    );

    const title = <h1 className="cursor-pointer text-lg font-semibold tracking-tight" onClick={() => navigate("/")}>FoodSmash</h1>;

    return (
        <nav className="sticky top-0 z-40 h-16 w-full shrink-0 border-b border-black bg-white px-4 py-3 sm:h-screen sm:w-48 sm:border-b-0 sm:border-r sm:py-6 md:w-64">
            {/* Mobile: brand + burger */}
            <div className="flex h-full items-center justify-between sm:hidden">
                {title}
                <IconButton onClick={() => setDrawerOpen(true)} variant="primary-outline" aria-label="Open menu">
                    <Menu size={16} />
                </IconButton>
            </div>

            {/* Desktop: brand + nav items */}
            <div className="hidden sm:block">
                <Stack space="lg">
                    {title}
                    {navItems}
                </Stack>
            </div>

            {/* Mobile drawer */}
            <MobileNavDrawer open={drawerOpen} onClose={() => setDrawerOpen(false)}>
                {navItems}
            </MobileNavDrawer>
        </nav>
    );
}