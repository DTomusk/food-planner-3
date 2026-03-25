import { Alert, Button, PageTitle, Spinner } from "@/components";
import Container from "@/components/layout/Container";
import Inline from "@/components/layout/Inline";
import Stack from "@/components/layout/Stack";
import SectionTitle from "@/components/ui/SectionTitle";
import { RecipeList, useRecipes } from "@/features/recipes";
import { Page } from "@/layout";
import { extractErrorMessage } from "@/lib/errors";
import { useSearchParams } from "react-router-dom";

export default function RecipeSearchPage() {
    const [searchParams] = useSearchParams();
    const query = searchParams.get("q")?.trim() ?? "";
    const hasQuery = query.length > 0;

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
                {!hasQuery && <Alert message="Enter a search query to see recipe results." closable={false} />}
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