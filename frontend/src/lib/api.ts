import { env } from "@/utils/constants"

export interface Image {
  id: string
  original_image_url: string
  object_storage_image_key: string
  transformed_image_key: string
  mime_type: string
  checksum: string
  status: string
  error_message?: string
  transformations: TransformationRequest[]
  created_at: string
  updated_at: string
}

export interface ImageMetadata {
  id: string
  original_image_url: string
  object_storage_image_key: string
  transformed_image_key: string
  mime_type: string
  status: string
  checksum: string
  error_message?: string
  transformation_count: number
  updated_at: string
  created_at: string
}

export type TransformationRequest =
  | { name: "resize"; config: { width: number; height: number } }
  | { name: "grayscale"; config: Record<string, never> }
  | { name: "trim"; config: { threshold: number } }
  | { name: "blur"; config: { sigma: number } }
  | { name: "rotate"; config: { angle: 90 | 180 | 270 } }

export interface CreateImageRequest {
  image_url: string
  transformations: TransformationRequest[]
}

export interface UpdateImageRequest {
  id: string
  transformations: TransformationRequest[]
}

export interface ListImagesResponse {
  page: number
  limit: number
  count: number
  data: Image[]
}

const API_BASE = env.VITE_API_URL

export const api = {
  async listImages(page = 1, limit = 10): Promise<ListImagesResponse> {
    const res = await fetch(
      `${API_BASE}/v1/images/?page=${page}&limit=${limit}`,
    )
    if (!res.ok) throw new Error("Failed to fetch images")
    return res.json()
  },

  async getImage(id: string): Promise<Image> {
    const res = await fetch(`${API_BASE}/v1/images/${id}`)
    if (!res.ok) throw new Error("Failed to fetch image")
    return res.json()
  },

  async getImageMetadata(id: string): Promise<ImageMetadata> {
    const res = await fetch(`${API_BASE}/v1/images/${id}/metadata`)
    if (!res.ok) throw new Error("Failed to fetch image metadata")
    return res.json()
  },

  async createImage(data: CreateImageRequest): Promise<{ id: string }> {
    const res = await fetch(`${API_BASE}/v1/images/`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    })
    if (!res.ok) throw new Error("Failed to create image")
    return res.json()
  },

  async updateImage(data: UpdateImageRequest): Promise<{ message: string }> {
    const res = await fetch(`${API_BASE}/v1/images/`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    })
    if (!res.ok) throw new Error("Failed to update image")
    return res.json()
  },

  async deleteImage(id: string): Promise<{ message: string }> {
    const res = await fetch(`${API_BASE}/v1/images/${id}`, {
      method: "DELETE",
    })
    if (!res.ok) throw new Error("Failed to delete image")
    return res.json()
  },

  streamAllImages(onMessage: (image: Image) => void): EventSource {
    const eventSource = new EventSource(`${API_BASE}/v1/images/sse`)
    eventSource.onmessage = (event) => {
      const image = JSON.parse(event.data)
      onMessage(image)
    }
    return eventSource
  },

  streamImage(id: string, onMessage: (image: Image) => void): EventSource {
    const eventSource = new EventSource(`${API_BASE}/v1/images/${id}/sse`)
    eventSource.onmessage = (event) => {
      const image = JSON.parse(event.data)
      onMessage(image)
    }
    return eventSource
  },
}
