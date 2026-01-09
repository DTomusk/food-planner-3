import { useNavigate } from "react-router-dom";

export default function BackLink() {
    const navgiate = useNavigate();
    return (
        <button
            className="text-blue-500 hover:underline mb-4"
            onClick={() => navgiate(-1)}
        >
            &larr; Back
        </button>
    );
}