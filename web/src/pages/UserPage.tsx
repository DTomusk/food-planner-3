import { Alert, PageTitle, Spinner } from "@/components";
import Container from "@/components/layout/Container";
import Stack from "@/components/layout/Stack";
import { useUser } from "@/features/users/hooks/useUser";
import { Page } from "@/layout";
import { extractErrorMessage } from "@/lib/errors";
import { useParams } from "react-router-dom";

export default function UserPage() {
    const { id } = useParams<{ id: string }>();
    const { data: user, isLoading, error } = useUser(id!);

    return (
        <Page>
            {isLoading && <Spinner/>}
            {error && <Alert message={extractErrorMessage(error)} closable />}
            <PageTitle text={user ? user.user?.username! : "User page"} />
            <Container size="xl">
                <Stack space="xl">
                    User page for user with id: {id}
                </Stack>
            </Container>
        </Page>
    );
}