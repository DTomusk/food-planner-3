interface SharedByProps {
    userName: string;
}

// TODO: add user id and link to user page
// TODO: maybe move into user feature? Or recipe...
export default function SharedBy({ userName }: SharedByProps) {
    return (
        <p className="text-center text-sm text-gray-500">
            Shared by {userName}
        </p>
    );
}