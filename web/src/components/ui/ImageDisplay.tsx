import { useState } from "react";

type ImageDisplayProps = {
    imageUrl: string | null;
    altText?: string;
};

export default function ImageDisplay({ imageUrl, altText }: ImageDisplayProps) {
    const [loading, setLoading] = useState(true);
    const [imgSrc, setImgSrc] = useState(imageUrl);

    const handleError = () => {
        setLoading(false);
        setImgSrc("https://placehold.co/600x400?text=Image+Not+Found");
    };

    return (
        <>
            {loading && <div className="w-full h-64 bg-gray-200 animate-pulse" />}
            <img 
                src={imgSrc ?? undefined} 
                alt={altText} 
                onError={handleError}
                onLoad={() => setLoading(false)}
                className={`w-full h-40 object-cover transition-opacity duration-300 overflow-hidden
                ${
                    loading ? "opacity-0" : "opacity-100"
                }`}
            />
        </>
    );
}