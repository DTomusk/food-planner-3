import type { User } from "@/features/recipes/types";
import Link from "./ui/Link";
import { useNavigate } from "react-router-dom";
import Inline from "./layout/Inline";

interface SharedByProps {
    user: User;
}

// TODO: maybe move into user feature?
export default function SharedBy({ user }: SharedByProps) {
    const navigate = useNavigate();
    return (
        <Inline align="center">
            <Link onClick={() => navigate(`/users/${user.id}`)}>
                Shared by {user.username}
            </Link>
        </Inline>
    );
}