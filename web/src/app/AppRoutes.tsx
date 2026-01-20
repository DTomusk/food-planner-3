import { Outlet, Route, Routes } from "react-router-dom";
import HomePage from "../pages/HomePage";
import RecipePage from "../pages/RecipePage";
import SignInPage from "@/pages/SignInPage";
import SignUpPage from "@/pages/SignUpPage";
import NotFound from "@/pages/NotFound";
import RecipeListingPage from "@/pages/RecipeListingPage";
import RecipeCreatePage from "@/pages/RecipeCreatePage";
import ProtectedLayout from "./ProtectedLayout";

function AuthLayout() {
  return (
    <Outlet />
  );
}

function RecipeLayout() {
  return (
    <Outlet />
  );
}

export function AppRoutes() {
  return (
    <Routes>
        <Route path="/" element={<HomePage/>}/>
        <Route path="/recipe" element={<RecipeLayout />}>
          <Route index element={<RecipeListingPage />} />
          <Route element={<ProtectedLayout />}>
            <Route path="create" element={<RecipeCreatePage />} />
          </Route>
          <Route path=":id" element={<RecipePage />} />
        </Route>
        <Route path="/auth" element={<AuthLayout />}>
          <Route path="signin" element={<SignInPage />} />
          <Route path="signup" element={<SignUpPage />} />
        </Route>
        <Route path="*" element={<NotFound />} />
    </Routes>
  );
}