import { createBrowserRouter, Outlet, type LoaderFunctionArgs } from "react-router-dom";
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
import type { GetRecipeQuery, GetRecipeVersionQuery, GetUserQuery } from "@/lib";
import { queryClient } from "./queryClient";
import { recipeQueryOptions } from "@/features/recipes/hooks/useRecipe";
import { recipeVersionQueryOptions } from "@/features/recipes/hooks/useRecipeVersion";
import { userQueryOptions } from "@/features/users/hooks/useUser";

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

async function loadRecipe({ params }: LoaderFunctionArgs) {
  const recipeId = params.id;

  if (!recipeId) {
    throw new Response("Recipe not found", { status: 404 });
  }

  return queryClient.ensureQueryData(recipeQueryOptions(recipeId));
}

async function loadRecipeVersion({ params }: LoaderFunctionArgs) {
  const recipeId = params.id;
  const version = params.version ? Number.parseInt(params.version, 10) : Number.NaN;

  if (!recipeId || Number.isNaN(version)) {
    throw new Response("Recipe version not found", { status: 404 });
  }

  return queryClient.ensureQueryData(recipeVersionQueryOptions(recipeId, version));
}

async function loadUser({ params }: LoaderFunctionArgs) {
  const userId = params.id;

  if (!userId) {
    throw new Response("User not found", { status: 404 });
  }

  return queryClient.ensureQueryData(userQueryOptions(userId));
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
            ],
          },
          {
            path: ":id",
            loader: loadRecipe,
            handle: {
              crumb: (data?: unknown) => {
                const recipeData = data as GetRecipeQuery | undefined;

                return recipeData?.recipe?.currentVersion.name ?? "Recipe";
              },
            },
            children: [
              {
                index: true,
                element: <RecipePage />,
              },
              {
                path: "versions/:version",
                loader: loadRecipeVersion,
                handle: {
                  crumb: (data?: unknown) => {
                    const recipeVersionData = data as GetRecipeVersionQuery | undefined;
                    const recipeName = recipeVersionData?.recipe?.version?.name;
                    const version = recipeVersionData?.recipe?.version?.version;

                    if (recipeName && version) {
                      return `v${version} ${recipeName}`;
                    }

                    return "Version";
                  },
                },
                element: <RecipeVersionPage />,
              },
              {
                element: <ProtectedLayout />,
                children: [
                  {
                    path: "edit",
                    element: <RecipeUpdatePage />,
                    handle: { crumb: () => "Edit" },
                  },
                ],
              },
            ],
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
        loader: loadUser,
        handle: {
          crumb: (data?: unknown) => {
            const userData = data as GetUserQuery | undefined;

            return userData?.user?.username ?? "User";
          },
        },
        element: <UserPage />,
      },
      {
        path: "*",
        element: <NotFound />,
      },
    ],
  }
])