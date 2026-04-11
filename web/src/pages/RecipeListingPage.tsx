import { Alert, Button, PageTitle, Spinner } from "@/components";
import Inline from "@/components/layout/Inline";
import Stack from "@/components/layout/Stack";
import SearchBar from "@/components/ui/SearchBar";
import { useAuth } from "@/features/auth/hooks/useAuth";
import { RecipeList, useRecipes } from "@/features/recipes";
import FilterChip from "@/features/recipes/components/FilterChip";
import RecipeFilters from "@/features/recipes/components/RecipeFilters";
import { Page } from "@/layout";
import { extractErrorMessage } from "@/lib/errors";
import { commonStrings } from "@/lib/strings";
import { Plus } from "lucide-react";
import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";

type DietLevelParam = "all" | "0" | "1";

function normalizeDietLevelParam(value: string | null): DietLevelParam {
    if (value === "0" || value === "1") {
        return value;
    }
    return "all";
}

function dietLevelLabel(dietLevel: DietLevelParam): string | null {
    if (dietLevel === "0") {
        return "Vegan";
    }
    if (dietLevel === "1") {
        return "Vegetarian";
    }
    return null;
}

export default function RecipeListingPage() {
    const [searchParams, setSearchParams] = useSearchParams();
    const query = searchParams.get("q")?.trim() ?? "";
    const hasQuery = query.length > 0;
    const [draftQuery, setDraftQuery] = useState(query);

    const dietLevel = normalizeDietLevelParam(searchParams.get("dietLevel"));

    // Keep input value in sync when URL changes via navigation/back-forward.
    useEffect(() => {
        setDraftQuery(query);
    }, [query]);

    const handleSubmitSearch = () => {
        const nextQuery = draftQuery.trim();
        const nextParams = new URLSearchParams(searchParams);

        if (nextQuery.length > 0) {
            nextParams.set("q", nextQuery);
        } else {
            nextParams.delete("q");
        }

        setSearchParams(nextParams);
    };

    const handleDietLevelChange = (nextDietLevel: DietLevelParam) => {
        const nextParams = new URLSearchParams(searchParams);
        if (nextDietLevel === "all") {
            nextParams.delete("dietLevel");
        } else {
            nextParams.set("dietLevel", nextDietLevel);
        }
        setSearchParams(nextParams);
    };

    const clearSearchFilter = () => {
        const nextParams = new URLSearchParams(searchParams);
        nextParams.delete("q");
        setSearchParams(nextParams);
    };

    const clearDietFilter = () => {
        const nextParams = new URLSearchParams(searchParams);
        nextParams.delete("dietLevel");
        setSearchParams(nextParams);
    };

    const clearAllFilters = () => {
        const nextParams = new URLSearchParams(searchParams);
        nextParams.delete("q");
        nextParams.delete("dietLevel");
        setSearchParams(nextParams);
    };

    const animalProductLevel = dietLevel === "all" ? undefined : Number(dietLevel);
    const activeDietLabel = dietLevelLabel(dietLevel);
    const hasActiveFilters = hasQuery || activeDietLabel !== null;

    const { data, isLoading, error, hasNextPage, fetchNextPage, isFetchingNextPage } = useRecipes({
        first: 20,
        query,
        animalProductLevel,
    });
    const { isAuthenticated } = useAuth();
    
    const navigate = useNavigate();
    return (
        <Page>
            <PageTitle text="Recipes" />
            <div className="mx-auto px-4 sm:px-6 lg:px-8">
                <Stack space="lg">
                {!isAuthenticated && <Alert type="info" message="Please sign in to add a new recipe." />}
                {isAuthenticated && 
                <Inline>
                    <Stack space="sm">
                        <Button onClick={() => navigate("/recipes/create")} 
                        aria-label="Add new recipe" variant="primary">
                            <Inline>
                            <Plus /> {commonStrings.recipe.add}
                            </Inline>
                        </Button>
                    </Stack>
                </Inline>}
                {error && <Alert message={extractErrorMessage(error)} closable />}
                <Inline align="start" justify="start" wrap className="w-full" gap="lg">
                    {/* Filters */}
                    <RecipeFilters dietLevel={dietLevel} onDietLevelChange={handleDietLevelChange} />
                    {/* Search bar and list */}
                    <Stack className="min-w-0 flex-1" space="lg">
                        <SearchBar
                            value={draftQuery}
                            onChange={setDraftQuery}
                            onSubmit={handleSubmitSearch}
                        />
                        {hasActiveFilters && (
                            <Inline justify="start" align="center" className="w-full rounded-md border border-black bg-gray-50 p-2">
                                <span className="text-xs font-semibold uppercase tracking-wide text-gray-600">Active filters</span>
                                {hasQuery && (
                                    <FilterChip label={`Search: "${query}"`} onClear={clearSearchFilter} />
                                )}
                                {activeDietLabel && (
                                    <FilterChip label={`Diet: ${activeDietLabel}`} onClear={clearDietFilter} />
                                )}
                                <button
                                    type="button"
                                    onClick={clearAllFilters}
                                    className="text-xs font-semibold text-gray-700 underline underline-offset-2 hover:text-gray-900 hover:cursor-pointer"
                                >
                                    Clear all
                                </button>
                            </Inline>
                        )}
                        {isLoading && <Spinner/>}
                        {data && (
                        <>
                            <RecipeList recipes={data.recipes} />
                            {hasNextPage && (
                                <Inline>
                                    <Button
                                        onClick={() => {
                                            void fetchNextPage();
                                        }}
                                        variant="secondary"
                                        loading={isFetchingNextPage}
                                    >
                                        Load more
                                    </Button>
                                </Inline>
                            )}
                        </>
                        )}
                    </Stack>
                </Inline>
                </Stack>
            </div>
        </Page>
    );
}