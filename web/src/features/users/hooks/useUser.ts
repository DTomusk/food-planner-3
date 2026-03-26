import { requireRecipeConnection, toRecipeSummaryList } from "@/features/recipes";
import type { RecipeSummary, User } from "@/features/recipes/types";
import { GetUserDocument, graphqlRequest } from "@/lib";
import type { GetUserQuery } from "@/lib/graphql.generated";
import { useQuery } from "@tanstack/react-query";
import type { ClientError } from "graphql-request";

type UserWithRecipes = {
    user: User;
    recipes: RecipeSummary[];
};

// TODO: this will probably be split out into multiple hooks, one to get user details, another to get some recipes, idk
export function useUser(id: string) {
    // TODO: define a type here that we don't pull from recipes
    return useQuery<GetUserQuery, ClientError, UserWithRecipes>({
        queryKey: ["user", id],
        queryFn: () => graphqlRequest(GetUserDocument, { id }),
        enabled: Boolean(id),
        select: (data) => {
            if (!data.user) {
                throw new Error("User not found");
            }

            return {
                user: {
                    id: data.user.id,
                    username: data.user.username,
                },
                recipes: toRecipeSummaryList(requireRecipeConnection(data.user.recipes)),
            };
        },
    });
}