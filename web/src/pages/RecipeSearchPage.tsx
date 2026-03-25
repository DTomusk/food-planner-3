import { Alert, Button, PageTitle, Spinner } from "@/components";
import Container from "@/components/layout/Container";
import Inline from "@/components/layout/Inline";
import Stack from "@/components/layout/Stack";
import SearchBar from "@/components/ui/SearchBar";
import SectionTitle from "@/components/ui/SectionTitle";
import { RecipeList, useRecipes } from "@/features/recipes";
import { Page } from "@/layout";
import { extractErrorMessage } from "@/lib/errors";
import { useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";

export default function RecipeSearchPage() {
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
        enabled: hasQuery,
    });
    
    return (
        <Page>
            <PageTitle text="Recipes" />
            <Container size="md">
                <Stack space="lg">
                <SearchBar
                    value={draftQuery}
                    onChange={setDraftQuery}
                    onSubmit={handleSubmitSearch}
                />
                {isLoading && <Spinner/>}
                {error && <Alert message={extractErrorMessage(error)} closable />}
                {hasQuery && data && (
                    <>
                        <SectionTitle text="Search results" />
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