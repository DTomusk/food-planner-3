import clsx from "clsx";
import { useEffect, useState } from "react";

const fallbackImageUrl = "https://placehold.co/600x400?text=Image+Not+Found";

type ImageDisplayProps = {
    imageUrl: string | null;
    altText?: string;
    containerClassName?: string;
    imageClassName?: string;
};

export default function ImageDisplay({
    imageUrl,
    altText,
    containerClassName,
    imageClassName,
}: ImageDisplayProps) {
    const [loading, setLoading] = useState(true);
    const [imgSrc, setImgSrc] = useState(imageUrl ?? fallbackImageUrl);

    useEffect(() => {
        setImgSrc(imageUrl ?? fallbackImageUrl);
        setLoading(true);
    }, [imageUrl]);

    const handleError = () => {
        if (imgSrc === fallbackImageUrl) {
            setLoading(false);
            return;
        }

        setImgSrc(fallbackImageUrl);
        setLoading(true);
    };

    return (
        <div className={clsx("relative overflow-hidden bg-gray-200", containerClassName)}>
            {loading && <div aria-hidden="true" className="absolute inset-0 animate-pulse bg-gray-200" />}
            <img
                src={imgSrc}
                alt={altText}
                onError={handleError}
                onLoad={() => setLoading(false)}
                className={clsx(
                    "h-full w-full object-cover transition-opacity duration-300",
                    loading ? "opacity-0" : "opacity-100",
                    imageClassName,
                )}
            />
        </div>
    );
}