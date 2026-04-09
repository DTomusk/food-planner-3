import { createBrowserRouter, Outlet, } from "react-router-dom";
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
import { AppLayout } from "@/layout";

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

export const router = createBrowserRouter([
  {
    path: "/",
    element: <AppLayout />,
    handle: { crumb: () => "Home" },
    children: [
      {
        index: true,
        element: <HomePage />,
      },
      {
        path: "recipes",
        element: <RecipeLayout />,
        handle: { crumb: () => "Recipes" },
        children: [
          {
            index: true,
            element: <RecipeListingPage />,
          },
          {
            element: <ProtectedLayout />,
            children: [
              {
                path: "create",
                element: <RecipeCreatePage />,
                handle: { crumb: () => "Create" },
              },
              {
                path: ":id/edit",
                element: <RecipeUpdatePage />,
                handle: { crumb: () => "Edit" },
              },
            ],
          },
          {
            path: ":id/versions/:version",
            element: <RecipeVersionPage />,
          },
          {
            path: ":id",
            element: <RecipePage />,
          },
        ],
      },
      {
        path: "auth",
        element: <AuthLayout />,
        children: [
          {
            path: "signin",
            element: <SignInPage />,
            handle : { crumb: () => "Sign In" },
          },
          {
            path: "signup",
            element: <SignUpPage />,
            handle : { crumb: () => "Sign Up" },
          },
        ],
      },
      {
        element: <ProtectedLayout />,
        children: [
          {
            path: "me/recipes",
            element: <MyRecipesPage />,
            handle: { crumb: () => "My Recipes" },
          },
        ],
      },
      {
        path: "users/:id",
        element: <UserPage />,
      },
      {
        path: "*",
        element: <NotFound />,
      },
    ],
  }
])