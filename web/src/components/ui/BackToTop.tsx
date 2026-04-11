import { ArrowUp } from "lucide-react";
import IconButton from "./IconButton";

type BackToTopProps = {
    showScrollTopButton: boolean;
    scrollToTop: () => void;
};

export default function BackToTop({ showScrollTopButton, scrollToTop }: BackToTopProps) {
    return (
        <IconButton
            variant="secondary"
            shape="circle"
            onClick={scrollToTop}
            aria-label="Back to top"
            className={[
                "fixed bottom-6 right-6 z-50 shadow-md",
                "transition-all duration-200",
                showScrollTopButton ? "translate-y-0 opacity-100" : "pointer-events-none translate-y-2 opacity-0",
            ].join(" ")}
        >
            <ArrowUp size={18} />
        </IconButton>
    )
}