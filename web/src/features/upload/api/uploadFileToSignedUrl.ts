export type UploadFileToSignedUrlInput = {
  uploadUrl: string;
  file: File;
  fileType?: string;
};

export async function uploadFileToSignedUrl({
  uploadUrl,
  file,
  fileType,
}: UploadFileToSignedUrlInput): Promise<void> {
  const response = await fetch(uploadUrl, {
    method: "PUT",
    headers: {
      "Content-Type": fileType ?? file.type ?? "application/octet-stream",
    },
    body: file,
  });

  if (!response.ok) {
    throw new Error(`Image upload failed (${response.status})`);
  }
}
