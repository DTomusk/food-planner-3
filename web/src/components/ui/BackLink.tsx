import { useNavigate } from "react-router-dom";

interface BackLinkProps {
    to: string;
}

export default function BackLink({ to }: BackLinkProps) {
    const navigate = useNavigate();
    return (
        <button
            className="text-blue-500 hover:underline mb-4"
            onClick={() => navigate(to)}
        >
            &larr; Back
        </button>
    );
}