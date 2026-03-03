import type { User } from "@/features/recipes/types";
import Link from "./ui/Link";
import { useNavigate } from "react-router-dom";

interface SharedByProps {
    user: User;
}

// TODO: add user id and link to user page
// TODO: maybe move into user feature? Or recipe...
export default function SharedBy({ user }: SharedByProps) {
    const navigate = useNavigate();
    return (
        <Link onClick={() => navigate(`/users/${user.id}`)} text={`Shared by ${user.username}`} />
    );
}