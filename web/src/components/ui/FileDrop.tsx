import { UploadCloud } from "lucide-react";
import { useDropzone, type Accept } from "react-dropzone";
import Button from "./Button";
import Text from "./Text";
import Stack from "../layout/Stack";
import ImageDisplay from "./ImageDisplay";

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
            className={`w-full border-2 border-dashed rounded-md flex flex-col items-center justify-center cursor-pointer transition-colors 
            ${isDragActive ? "border-primary-600 bg-primary-50" : "border-gray-300 hover:border-primary-500"}`}
            onClick={open}
        >
            <input {...getInputProps()} />
            {!value ? (
            <Stack space="md" className="py-4 px-4 items-center">
                <UploadCloud className="w-10 h-10 text-gray-400" />
                <Text variant="muted">{label}</Text>

                <Button variant="primaryOutline" onClick={open}>
                    Choose file
                </Button>
            </Stack>) : (
            <Preview file={value} onRemove={() => onChange(null)} />
            )}
        </div>
    )
}

function Preview({ file, onRemove } : { file: File; onRemove: () => void }) {
  const preview = URL.createObjectURL(file);

  return (
    <Stack space="md" className="py-4 px-4 items-center">
        <ImageDisplay imageUrl={preview} 
            altText={file.name}
            containerClassName="w-24 h-24 rounded-md overflow-hidden"
            imageClassName="object-contain"
        />
        <Text>{file.name}</Text>

        <Button
            variant="dangerOutline"
            onClick={onRemove}
        >
        Remove
        </Button>
    </Stack>
  );
}