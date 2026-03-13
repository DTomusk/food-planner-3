import { useLocation, useNavigate } from "react-router-dom";
import Link from "./Link";

export default function BackLink() {
    const navigate = useNavigate();
  const location = useLocation();

  const getFallback = () => {
    const parts = location.pathname.split("/").filter(Boolean);

    if (parts.length <= 1) {
      return "/";
    }

    parts.pop();
    return "/" + parts.join("/");
  };

  const handleBack = (e: React.MouseEvent) => {
    e.preventDefault();

    if (window.history.length > 1) {
      navigate(-1);
    } else {
      navigate(getFallback());
    }
  };
    return (
        <Link
            onClick={handleBack}
            color="primary"
        >
            &larr; Back
        </Link>
    );
}