import { Alert, PageTitle, Spinner } from "@/components";
import Container from "@/components/layout/Container";
import Stack from "@/components/layout/Stack";
import SectionTitle from "@/components/ui/SectionTitle";
import { RecipeList } from "@/features/recipes";
import { useUser } from "@/features/users/hooks/useUser";
import { Page } from "@/layout";
import { extractErrorMessage } from "@/lib/errors";
import { useParams } from "react-router-dom";

export default function UserPage() {
    const { id } = useParams<{ id: string }>();
    const { data: { user, recipes } = {}, isLoading, error } = useUser(id!);

    return (
        <Page>
            <PageTitle text={user ? user.username : ""} />
            <Container size="xl">
                <Stack space="xl">
                {isLoading && <Spinner/>}
                {error && <Alert message={extractErrorMessage(error)} closable />}
                {recipes && recipes.length > 0 && (
                    <>
                        <SectionTitle text="Top recipes" />
                        <RecipeList recipes={recipes} />
                    </>
                )}
                </Stack>
            </Container>
        </Page>
    );
}