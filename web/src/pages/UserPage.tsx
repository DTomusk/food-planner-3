import Container from "@/components/layout/Container";
import Stack from "@/components/layout/Stack";
import { Page } from "@/layout";
import { useParams } from "react-router-dom";

export default function UserPage() {
    const { id } = useParams<{ id: string }>();
    return (
        <Page>
            <Container size="xl">
                <Stack space="xl">
                    User page for user with id: {id}
                </Stack>
            </Container>
        </Page>
    );
}