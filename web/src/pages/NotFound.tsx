import { PageTitle, Text } from "@/components";
import { Page } from "@/layout";

export default function NotFound() {
    return (
        <Page>
            <PageTitle text="404 - Page Not Found" />
            <Text align="center">The page you are looking for does not exist.</Text>
        </Page>
    );
}