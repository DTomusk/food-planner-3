import { queryClient } from "@/app/queryClient";
import { clearToken, getToken, setToken } from "@/lib/auth/token";
import { setUnauthenticatedHandler } from "@/lib/auth/unauthenticated";
import { createContext, useEffect, useState, type ReactNode } from "react";

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
        queryClient.clear();
        setIsAuthenticated(false);
    };

    useEffect(() => {
        setUnauthenticatedHandler(signOut);

        return () => {
            setUnauthenticatedHandler(null);
        };
    }, [signOut]);

    return (
        <AuthContext.Provider value={{ isAuthenticated, signIn, signOut }}>
            {children}
        </AuthContext.Provider>
    );  
}