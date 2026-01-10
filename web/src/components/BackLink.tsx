import { useNavigate } from "react-router-dom";

// TODO: consider injecting route 
// Right now, if you navigate to a page directly,
// the back button takes you back to the last page you were on, not up the hierarchy
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