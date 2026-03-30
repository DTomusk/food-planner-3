import { UploadCloud } from "lucide-react";
import { useDropzone, type Accept } from "react-dropzone";
import Button from "./Button";
import Text from "./Text";

type FileDropProps = {
    value: File | null;
    onChange: (file: File | null) => void;
    accept?: Accept;
    maxSize?: number;
    label?: string;
};

export default function FileDrop({ 
    value, 
    onChange, 
    label = "Select a file",
    accept,
    maxSize
}: FileDropProps) {
    const { getRootProps, getInputProps, open, isDragActive } = useDropzone({
        accept,
        maxSize,
        multiple: false,
        noClick: true,
        onDrop: (acceptedFiles) => {
            if (acceptedFiles[0]) {
                onChange(acceptedFiles[0]);
            }
        },
    });

    return (
        <div
            {...getRootProps()}
            className={`w-full h-32 border-2 border-dashed rounded-md flex flex-col items-center justify-center cursor-pointer transition-colors ${isDragActive ? "border-primary-600 bg-primary-50" : "border-gray-300 hover:border-gray-500"}`}
            onClick={open}
        >
            <input {...getInputProps()} />
            {!value && (
        <div className="flex flex-col items-center gap-3">
            <UploadCloud className="w-10 h-10 text-gray-400" />
            <Text variant="muted">{label}</Text>

            <Button variant="primaryOutline" onClick={open}>
                Choose file
            </Button>
        </div>)}
        </div>
    )
}