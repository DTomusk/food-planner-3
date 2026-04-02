import { useMutation } from "@tanstack/react-query";
import {
  uploadFileToSignedUrl,
  type UploadFileToSignedUrlInput,
} from "../api/uploadFileToSignedUrl";

export function useUploadFileToSignedUrl() {
  return useMutation<void, Error, UploadFileToSignedUrlInput, unknown>({
    mutationFn: uploadFileToSignedUrl,
  });
}
