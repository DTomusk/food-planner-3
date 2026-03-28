import { Alert, Button, PageTitle, Spinner } from "@/components";
import Container from "@/components/layout/Container";
import Inline from "@/components/layout/Inline";
import Stack from "@/components/layout/Stack";
import SearchBar from "@/components/ui/SearchBar";
import SectionTitle from "@/components/ui/SectionTitle";
import { useAuth } from "@/features/auth/hooks/useAuth";
import { RecipeList, useRecipes } from "@/features/recipes";
import { Page } from "@/layout";
import { extractErrorMessage } from "@/lib/errors";
import { commonStrings } from "@/lib/strings";
import { Plus } from "lucide-react";
import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";

export default function RecipeListingPage() {
    const [searchParams, setSearchParams] = useSearchParams();
    const query = searchParams.get("q")?.trim() ?? "";
    const hasQuery = query.length > 0;
    const [draftQuery, setDraftQuery] = useState(query);

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

    const { data, isLoading, error, hasNextPage, fetchNextPage, isFetchingNextPage } = useRecipes({
        first: 20,
        query,
    });
    const { isAuthenticated } = useAuth();
    
    const navigate = useNavigate();
    return (
        <Page>
            <PageTitle text="Recipes" />
            <Container size="md">
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
                {isLoading && <Spinner/>}
                {error && <Alert message={extractErrorMessage(error)} closable />}
                <SearchBar
                    value={draftQuery}
                    onChange={setDraftQuery}
                    onSubmit={handleSubmitSearch}
                />
                {data && (
                    <>
                        <SectionTitle text={hasQuery ? "Search results" : "Top recipes"} />
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
            </Container>
        </Page>
    );
}