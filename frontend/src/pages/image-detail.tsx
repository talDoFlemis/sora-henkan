import { useEffect, useState } from "react"
import { useParams, useNavigate } from "react-router-dom"
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import { api, type Image, type TransformationRequest, type ImageMetadata } from "@/lib/api"
import { ImageComparison } from "@/components/image-comparison"
import {
  ArrowLeft,
  Trash2,
  Edit,
  Save,
  X,
  Plus,
  Image as ImageIcon,
  Loader2,
  ExternalLink,
  Copy,
  Check,
  Home,
  Sparkles,
  Calendar,
  Hash,
  FileType,
  Maximize2,
  Palette,
  Scissors,
  Droplet,
  RotateCw,
  Database,
} from "lucide-react"
import { env } from "@/utils/constants"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Badge } from "@/components/ui/badge"

const transformationIcons: Record<string, React.ElementType> = {
  resize: Maximize2,
  grayscale: Palette,
  trim: Scissors,
  blur: Droplet,
  rotate: RotateCw,
}

const transformationColors: Record<string, string> = {
  resize: "from-indigo-500 to-purple-600",
  grayscale: "from-gray-500 to-gray-700",
  trim: "from-orange-500 to-red-500",
  blur: "from-cyan-500 to-blue-600",
  rotate: "from-pink-500 to-rose-600",
}

export function ImageDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [image, setImage] = useState<Image | null>(null)
  const [metadata, setMetadata] = useState<ImageMetadata | null>(null)
  const [loading, setLoading] = useState(true)
  const [editing, setEditing] = useState(false)
  const [transformations, setTransformations] = useState<
    TransformationRequest[]
  >([])
  const [updating, setUpdating] = useState(false)
  const [deleting, setDeleting] = useState(false)
  const [copiedId, setCopiedId] = useState(false)

  const loadImage = async () => {
    if (!id) return
    setLoading(true)
    try {
      const data = await api.getImage(id)
      setImage(data)
      setTransformations(data.transformations)

      try {
        const meta = await api.getImageMetadata(id)
        setMetadata(meta)
      } catch (err) {
        console.error("Failed to load metadata", err)
      }
    } catch (error) {
      console.error("Failed to load image", error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadImage()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id])

  useEffect(() => {
    if (!id) return
    const eventSource = api.streamImage(id, (updatedImage) => {
      setImage(updatedImage)
    })
    return () => eventSource.close()
  }, [id])

  const handleDelete = async () => {
    if (
      !id ||
      !confirm(
        "Are you sure you want to delete this image? This action cannot be undone.",
      )
    )
      return
    setDeleting(true)
    try {
      await api.deleteImage(id)
      navigate("/gallery")
    } catch {
      alert("Failed to delete image")
      setDeleting(false)
    }
  }

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!id) return
    setUpdating(true)
    try {
      await api.updateImage({ id, transformations })
      setEditing(false)
      await loadImage()
    } catch {
      alert("Failed to update image")
    } finally {
      setUpdating(false)
    }
  }

  const addTransformation = (type: TransformationRequest["name"]) => {
    const newTransform: TransformationRequest =
      type === "resize"
        ? { name: "resize", config: { width: 800, height: 600 } }
        : type === "grayscale"
          ? { name: "grayscale", config: {} }
          : type === "trim"
            ? { name: "trim", config: { threshold: 10 } }
            : type === "blur"
              ? { name: "blur", config: { sigma: 1.5 } }
              : { name: "rotate", config: { angle: 90 } }

    setTransformations([...transformations, newTransform])
  }

  const updateTransformation = (index: number, config: unknown) => {
    const updated = [...transformations]
    updated[index] = { ...updated[index], config } as TransformationRequest
    setTransformations(updated)
  }

  const removeTransformation = (index: number) => {
    setTransformations(transformations.filter((_, i) => i !== index))
  }

  const getImageUrl = (key: string) => `${env.AWS_BUCKET_ENDPOINT}/${key}`

  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text)
      setCopiedId(true)
      setTimeout(() => setCopiedId(false), 2000)
    } catch (error) {
      console.error("Failed to copy", error)
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case "processed":
      case "completed":
        return "bg-emerald-500"
      case "processing":
        return "bg-amber-500"
      case "pending":
        return "bg-blue-500"
      case "failed":
        return "bg-red-500"
      default:
        return "bg-gray-500"
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="w-20 h-20 rounded-2xl gradient-primary flex items-center justify-center mx-auto mb-6 animate-pulse">
            <Loader2 className="w-10 h-10 animate-spin text-white" />
          </div>
          <p className="text-gray-600 font-medium">Loading image details...</p>
        </div>
      </div>
    )
  }

  if (!image) {
    return (
      <div className="min-h-screen flex items-center justify-center p-6">
        <Card className="p-12 text-center max-w-md glass-dark rounded-3xl border-0 shadow-lg">
          <div className="w-20 h-20 rounded-2xl bg-gray-100 flex items-center justify-center mx-auto mb-6">
            <ImageIcon className="w-10 h-10 text-gray-400" />
          </div>
          <h2 className="text-2xl font-bold text-gray-800 mb-3">
            Image Not Found
          </h2>
          <p className="text-gray-600 mb-8">
            The image you're looking for doesn't exist or has been deleted.
          </p>
          <Button
            onClick={() => navigate("/gallery")}
            className="gradient-primary text-white rounded-xl px-6 py-3"
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back to Gallery
          </Button>
        </Card>
      </div>
    )
  }

  return (
    <div className="min-h-screen">
      {/* Header */}
      <header className="sticky top-0 z-50 glass-dark shadow-sm">
        <div className="w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center justify-between">
            <button
              onClick={() => navigate("/")}
              className="flex items-center gap-3 hover:opacity-80 transition-opacity"
            >
              <div className="w-10 h-10 rounded-xl gradient-primary flex items-center justify-center">
                <Sparkles className="w-5 h-5 text-white" />
              </div>
              <span className="text-xl font-bold gradient-text">
                Sora Henkan
              </span>
            </button>

            <div className="flex items-center gap-3">
              <Button
                variant="ghost"
                onClick={() => navigate("/")}
                className="text-gray-600 hover:text-gray-800"
              >
                <Home className="w-4 h-4 mr-2" />
                Home
              </Button>
              <Button
                variant="outline"
                onClick={() => navigate("/gallery")}
                className="rounded-xl"
              >
                <ArrowLeft className="w-4 h-4 mr-2" />
                Gallery
              </Button>
            </div>
          </div>
        </div>
      </header>

      <main className="w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Title & Actions */}
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8 animate-fade-in-up">
          <div>
            <div className="flex items-center gap-3 mb-2">
              <h1 className="text-3xl md:text-4xl font-bold text-gray-800">
                Image Details
              </h1>
              <Badge
                className={`${getStatusColor(image.status)} text-white border-0 rounded-lg px-3`}
              >
                {image.status}
              </Badge>
            </div>
            <div className="flex items-center gap-2 text-gray-500">
              <code className="text-sm bg-gray-100 px-2 py-1 rounded-lg">
                {image.id}
              </code>
              <Button
                variant="ghost"
                size="icon"
                className="h-8 w-8"
                onClick={() => copyToClipboard(image.id)}
              >
                {copiedId ? (
                  <Check className="w-4 h-4 text-emerald-500" />
                ) : (
                  <Copy className="w-4 h-4" />
                )}
              </Button>
            </div>
          </div>

          <div className="flex gap-3">
            <Button
              variant="outline"
              onClick={() => setEditing(!editing)}
              disabled={editing && updating}
              className="rounded-xl border-indigo-200 text-indigo-700 hover:bg-indigo-50"
            >
              {editing ? (
                <>
                  <X className="w-4 h-4 mr-2" />
                  Cancel
                </>
              ) : (
                <>
                  <Edit className="w-4 h-4 mr-2" />
                  Edit
                </>
              )}
            </Button>
            <Button
              variant="destructive"
              onClick={handleDelete}
              disabled={deleting}
              className="rounded-xl"
            >
              {deleting ? (
                <>
                  <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                  Deleting...
                </>
              ) : (
                <>
                  <Trash2 className="w-4 h-4 mr-2" />
                  Delete
                </>
              )}
            </Button>
          </div>
        </div>

        {/* Main Content */}
        <div className="grid lg:grid-cols-3 gap-8">
          {/* Image Comparison */}
          <div
            className="lg:col-span-2 animate-fade-in-up"
            style={{ animationDelay: "0.1s" }}
          >
            {editing ? (
              <Card className="p-6 glass-dark rounded-3xl border-0 shadow-lg">
                <h2 className="text-xl font-bold text-gray-800 mb-6">
                  Edit Transformations
                </h2>
                <form onSubmit={handleUpdate} className="space-y-6">
                  <div className="flex items-center justify-between">
                    <Label className="text-sm font-semibold text-gray-700">
                      Transformations
                    </Label>
                    <Select
                      onValueChange={(v) =>
                        addTransformation(v as TransformationRequest["name"])
                      }
                    >
                      <SelectTrigger className="w-[180px] rounded-xl">
                        <Plus className="w-4 h-4 mr-2" />
                        <SelectValue placeholder="Add transform" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="resize">Resize</SelectItem>
                        <SelectItem value="grayscale">Grayscale</SelectItem>
                        <SelectItem value="trim">Trim</SelectItem>
                        <SelectItem value="blur">Blur</SelectItem>
                        <SelectItem value="rotate">Rotate</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>

                  {transformations.length === 0 ? (
                    <div className="text-center py-8 text-gray-500">
                      <ImageIcon className="w-12 h-12 mx-auto mb-3 opacity-30" />
                      <p>No transformations added yet</p>
                    </div>
                  ) : (
                    <div className="space-y-3">
                      {transformations.map((transform, index) => {
                        const Icon =
                          transformationIcons[transform.name] || ImageIcon
                        const colorClass =
                          transformationColors[transform.name] ||
                          "from-gray-500 to-gray-600"

                        return (
                          <Card
                            key={index}
                            className="p-4 rounded-2xl border-gray-100"
                          >
                            <div className="flex items-start gap-4">
                              <div
                                className={`w-10 h-10 rounded-xl bg-gradient-to-br ${colorClass} flex items-center justify-center flex-shrink-0`}
                              >
                                <Icon className="w-5 h-5 text-white" />
                              </div>
                              <div className="flex-1 space-y-3">
                                <div className="flex items-center justify-between">
                                  <span className="font-semibold text-gray-800 capitalize">
                                    {transform.name}
                                  </span>
                                  <Button
                                    type="button"
                                    variant="ghost"
                                    size="icon"
                                    onClick={() => removeTransformation(index)}
                                    className="h-8 w-8 text-gray-400 hover:text-red-500"
                                  >
                                    <Trash2 className="w-4 h-4" />
                                  </Button>
                                </div>

                                {transform.name === "resize" && (
                                  <div className="grid grid-cols-2 gap-3">
                                    <div>
                                      <Label className="text-xs text-gray-500">
                                        Width
                                      </Label>
                                      <Input
                                        type="number"
                                        value={transform.config.width}
                                        onChange={(e) =>
                                          updateTransformation(index, {
                                            ...transform.config,
                                            width: +e.target.value,
                                          })
                                        }
                                        className="h-9 rounded-lg"
                                      />
                                    </div>
                                    <div>
                                      <Label className="text-xs text-gray-500">
                                        Height
                                      </Label>
                                      <Input
                                        type="number"
                                        value={transform.config.height}
                                        onChange={(e) =>
                                          updateTransformation(index, {
                                            ...transform.config,
                                            height: +e.target.value,
                                          })
                                        }
                                        className="h-9 rounded-lg"
                                      />
                                    </div>
                                  </div>
                                )}

                                {transform.name === "blur" && (
                                  <div>
                                    <Label className="text-xs text-gray-500">
                                      Sigma
                                    </Label>
                                    <Input
                                      type="number"
                                      step="0.1"
                                      value={transform.config.sigma}
                                      onChange={(e) =>
                                        updateTransformation(index, {
                                          sigma: +e.target.value,
                                        })
                                      }
                                      className="h-9 rounded-lg"
                                    />
                                  </div>
                                )}

                                {transform.name === "trim" && (
                                  <div>
                                    <Label className="text-xs text-gray-500">
                                      Threshold
                                    </Label>
                                    <Input
                                      type="number"
                                      value={transform.config.threshold}
                                      onChange={(e) =>
                                        updateTransformation(index, {
                                          threshold: +e.target.value,
                                        })
                                      }
                                      className="h-9 rounded-lg"
                                    />
                                  </div>
                                )}

                                {transform.name === "rotate" && (
                                  <div>
                                    <Label className="text-xs text-gray-500">
                                      Angle
                                    </Label>
                                    <Select
                                      value={transform.config.angle.toString()}
                                      onValueChange={(v) =>
                                        updateTransformation(index, {
                                          angle: +v,
                                        })
                                      }
                                    >
                                      <SelectTrigger className="h-9 rounded-lg">
                                        <SelectValue />
                                      </SelectTrigger>
                                      <SelectContent>
                                        <SelectItem value="90">90°</SelectItem>
                                        <SelectItem value="180">
                                          180°
                                        </SelectItem>
                                        <SelectItem value="270">
                                          270°
                                        </SelectItem>
                                      </SelectContent>
                                    </Select>
                                  </div>
                                )}
                              </div>
                            </div>
                          </Card>
                        )
                      })}
                    </div>
                  )}

                  <Button
                    type="submit"
                    disabled={updating || transformations.length === 0}
                    className="w-full h-12 gradient-primary text-white rounded-xl"
                  >
                    {updating ? (
                      <>
                        <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                        Saving...
                      </>
                    ) : (
                      <>
                        <Save className="w-4 h-4 mr-2" />
                        Save Changes
                      </>
                    )}
                  </Button>
                </form>
              </Card>
            ) : (
              <Card className="overflow-hidden glass-dark rounded-3xl border-0 shadow-lg">
                {image.status === "processed" && image.transformed_image_key ? (
                  <ImageComparison
                    beforeSrc={getImageUrl(image.object_storage_image_key)}
                    afterSrc={getImageUrl(image.transformed_image_key)}
                    alt="Original vs Transformed"
                  />
                ) : (
                  <div className="p-8">
                    <div className="aspect-video relative overflow-hidden rounded-2xl bg-gray-100">
                      <img
                        src={getImageUrl(image.object_storage_image_key)}
                        alt="Original"
                        className="w-full h-full object-contain"
                      />
                      {image.status === "processing" && (
                        <div className="absolute inset-0 bg-black/50 flex items-center justify-center backdrop-blur-sm">
                          <div className="text-center text-white">
                            <Loader2 className="w-12 h-12 animate-spin mx-auto mb-3" />
                            <p className="font-medium">Processing image...</p>
                            <p className="text-sm opacity-70">
                              This may take a moment
                            </p>
                          </div>
                        </div>
                      )}
                      {image.status === "failed" && (
                        <div className="absolute inset-0 bg-red-500/20 flex items-center justify-center backdrop-blur-sm">
                          <div className="text-center text-red-600 bg-white/90 p-6 rounded-2xl">
                            <X className="w-12 h-12 mx-auto mb-3" />
                            <p className="font-medium">Processing failed</p>
                            <p className="text-sm mt-2 text-gray-600">
                              {image.error_message || "Unknown error occurred"}
                            </p>
                          </div>
                        </div>
                      )}
                    </div>
                    {image.status === "pending" && (
                      <div className="mt-4 text-center text-gray-500">
                        <p className="text-sm">
                          Image is queued for processing
                        </p>
                      </div>
                    )}
                  </div>
                )}
              </Card>
            )}
          </div>

          {/* Sidebar Info */}
          <div
            className="space-y-6 animate-fade-in-up"
            style={{ animationDelay: "0.2s" }}
          >
            {/* Image Info */}
            <Card className="p-6 glass-dark rounded-3xl border-0 shadow-lg">
              <h3 className="font-bold text-gray-800 mb-4">
                Image Information
              </h3>
              <div className="space-y-4">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-xl bg-indigo-100 flex items-center justify-center">
                    <FileType className="w-5 h-5 text-indigo-600" />
                  </div>
                  <div>
                    <p className="text-xs text-gray-500">Mime Type</p>
                    <p className="font-medium text-gray-800">
                      {image.mime_type || "Unknown"}
                    </p>
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-xl bg-purple-100 flex items-center justify-center">
                    <Calendar className="w-5 h-5 text-purple-600" />
                  </div>
                  <div>
                    <p className="text-xs text-gray-500">Created</p>
                    <p className="font-medium text-gray-800">
                      {new Date(image.created_at).toLocaleDateString("en-US", {
                        month: "short",
                        day: "numeric",
                        year: "numeric",
                        hour: "2-digit",
                        minute: "2-digit",
                      })}
                    </p>
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-xl bg-pink-100 flex items-center justify-center">
                    <Hash className="w-5 h-5 text-pink-600" />
                  </div>
                  <div>
                    <p className="text-xs text-gray-500">Checksum</p>
                    <p className="font-medium text-gray-800 text-sm truncate max-w-[180px]">
                      {image.checksum || "Not available"}
                    </p>
                  </div>
                </div>
              </div>
            </Card>

            {/* Transformations */}
            <Card className="p-6 glass-dark rounded-3xl border-0 shadow-lg">
              <h3 className="font-bold text-gray-800 mb-4">
                Transformations ({image.transformations.length})
              </h3>
              {image.transformations.length === 0 ? (
                <p className="text-gray-500 text-sm">
                  No transformations applied
                </p>
              ) : (
                <div className="space-y-3">
                  {image.transformations.map((t, i) => {
                    const Icon = transformationIcons[t.name] || ImageIcon
                    const colorClass =
                      transformationColors[t.name] ||
                      "from-gray-500 to-gray-600"

                    return (
                      <div
                        key={i}
                        className="flex items-center gap-3 p-3 bg-gray-50 rounded-xl"
                      >
                        <div
                          className={`w-8 h-8 rounded-lg bg-gradient-to-br ${colorClass} flex items-center justify-center`}
                        >
                          <Icon className="w-4 h-4 text-white" />
                        </div>
                        <div>
                          <p className="font-medium text-gray-800 capitalize text-sm">
                            {t.name}
                          </p>
                          <p className="text-xs text-gray-500">
                            {t.name === "resize" &&
                              `${t.config.width}x${t.config.height}`}
                            {t.name === "blur" && `σ ${t.config.sigma}`}
                            {t.name === "rotate" && `${t.config.angle}°`}
                            {t.name === "trim" &&
                              `threshold ${t.config.threshold}`}
                            {t.name === "grayscale" && "Applied"}
                          </p>
                        </div>
                      </div>
                    )
                  })}
                </div>
              )}
            </Card>

            {/* Original URL */}
            <Card className="p-6 glass-dark rounded-3xl border-0 shadow-lg">
              <h3 className="font-bold text-gray-800 mb-4">Source URL</h3>
              <a
                href={image.original_image_url}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-2 text-sm text-indigo-600 hover:text-indigo-700 transition-colors break-all"
              >
                <ExternalLink className="w-4 h-4 flex-shrink-0" />
                <span className="truncate">{image.original_image_url}</span>
              </a>
            </Card>

            {/* Technical Metadata (DynamoDB) */}
            {metadata && (
              <Card className="p-6 glass-dark rounded-3xl border-0 shadow-lg">
                <div className="flex items-center gap-3 mb-4">
                  <div className="w-10 h-10 rounded-xl bg-green-100 flex items-center justify-center">
                    <Database className="w-5 h-5 text-green-600" />
                  </div>
                  <h3 className="font-bold text-gray-800">
                    Technical Metadata (DynamoDB)
                  </h3>
                </div>

                <div className="space-y-3 mb-4">
                  <div>
                    <p className="text-xs text-gray-500">Storage Key</p>
                    <p className="font-mono text-xs text-gray-800 break-all bg-gray-50 p-2 rounded-lg mt-1 border border-gray-100">
                      {metadata.object_storage_image_key}
                    </p>
                  </div>
                  {metadata.transformed_image_key && (
                    <div>
                      <p className="text-xs text-gray-500">Transformed Key</p>
                      <p className="font-mono text-xs text-gray-800 break-all bg-gray-50 p-2 rounded-lg mt-1 border border-gray-100">
                        {metadata.transformed_image_key}
                      </p>
                    </div>
                  )}
                  <div className="grid grid-cols-2 gap-3">
                    <div>
                      <p className="text-xs text-gray-500">Transformations</p>
                      <p className="font-medium text-gray-800">
                        {metadata.transformation_count}
                      </p>
                    </div>
                    <div>
                      <p className="text-xs text-gray-500">Checksum</p>
                      <p className="font-mono text-xs text-gray-800 truncate" title={metadata.checksum}>
                        {metadata.checksum.substring(0, 10)}...
                      </p>
                    </div>
                  </div>
                </div>

                <details className="group">
                  <summary className="text-xs font-semibold text-indigo-600 cursor-pointer hover:text-indigo-700 transition-colors flex items-center gap-1 select-none">
                    <span>View Raw JSON</span>
                  </summary>
                  <div className="mt-3 animate-fade-in-up">
                    <pre className="p-3 bg-gray-900 text-green-400 rounded-xl text-[10px] leading-relaxed overflow-x-auto font-mono custom-scrollbar">
                      {JSON.stringify(metadata, null, 2)}
                    </pre>
                  </div>
                </details>
              </Card>
            )}

          </div>
        </div>
      </main>
    </div>
  )
}
