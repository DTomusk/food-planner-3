import { PageTitle } from "@/components";
import { Page } from "@/layout";

export default function NotFound() {
    return (
        <Page>
            <PageTitle text="404 - Page Not Found" />
            <p>The page you are looking for does not exist.</p>
        </Page>
    );
}