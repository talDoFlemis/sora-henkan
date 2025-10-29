export const env = {
  VITE_API_URL: import.meta.env.VITE_API_URL ?? "http://localhost:42069",
  AWS_BUCKET_ENDPOINT:
    import.meta.env.VITE_AWS_BUCKET_ENDPOINT ?? "http://localhost:9000/images",
}
