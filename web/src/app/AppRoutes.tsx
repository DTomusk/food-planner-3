import { Outlet, Route, Routes } from "react-router-dom";
import HomePage from "../pages/HomePage";
import RecipePage from "../pages/RecipePage";
import SignInPage from "@/pages/SignInPage";
import SignUpPage from "@/pages/SignUpPage";
import NotFound from "@/pages/NotFound";
import RecipeListingPage from "@/pages/RecipeListingPage";
import RecipeCreatePage from "@/pages/RecipeCreatePage";
import RecipeVersionPage from "@/pages/RecipeVersionPage";
import ProtectedLayout from "./ProtectedLayout";
import MyRecipesPage from "@/pages/MyRecipesPage";
import UserPage from "@/pages/UserPage";
import RecipeUpdatePage from "@/pages/RecipeUpdatePage";
import RecipeSearchPage from "@/pages/RecipeSearchPage";

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
        <Route path="/recipes" element={<RecipeLayout />}>
          <Route index element={<RecipeListingPage />} />
          <Route element={<ProtectedLayout />}>
            <Route path="create" element={<RecipeCreatePage />} />
            <Route path=":id/edit" element={<RecipeUpdatePage />} />
          </Route>
          <Route path="search" element={<RecipeSearchPage />} />
          <Route path=":id/versions/:version" element={<RecipeVersionPage />} />
          <Route path=":id" element={<RecipePage />} />
        </Route>
        <Route path="/auth" element={<AuthLayout />}>
          <Route path="signin" element={<SignInPage />} />
          <Route path="signup" element={<SignUpPage />} />
        </Route>
        <Route element={<ProtectedLayout />}>
          <Route path="me/recipes" element={<MyRecipesPage />} />
        </Route>
        <Route path="/users/:id" element={<UserPage />} />
        <Route path="*" element={<NotFound />} />
    </Routes>
  );
}