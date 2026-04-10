import type { User } from "@/features/recipes/types";
import Link from "./ui/Link";
import { useNavigate } from "react-router-dom";

interface SharedByProps {
    user: User;
}

// TODO: maybe move into user feature?
export default function SharedBy({ user }: SharedByProps) {
    const navigate = useNavigate();
    return (
        <Link onClick={(e) => {
            e.preventDefault();
            e.stopPropagation();
            navigate(`/users/${user.id}`)}}>
            Shared by {user.username}
        </Link>
    );
}