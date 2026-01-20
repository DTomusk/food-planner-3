import { clearToken, getToken, setToken } from "@/lib/auth/token";
import { createContext, useState, type ReactNode } from "react";

type AuthContextValue = {
    isAuthenticated: boolean;
    signIn: (token: string) => void;
    signOut: () => void;
};

export const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
    const [isAuthenticated, setIsAuthenticated] = useState<boolean>(
        () => getToken() !== null
    );

    const signIn = (token: string) => {
        setToken(token);
        setIsAuthenticated(true);
    };

    const signOut = () => {
        clearToken();
        setIsAuthenticated(false);
    };

    return (
        <AuthContext.Provider value={{ isAuthenticated, signIn, signOut }}>
            {children}
        </AuthContext.Provider>
    );  
}