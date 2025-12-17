import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Card } from "@/components/ui/card"
import { api, type TransformationRequest } from "@/lib/api"
import {
  Trash2,
  Plus,
  Wand2,
  Link,
  Loader2,
  RotateCw,
  Maximize2,
  Droplet,
  Scissors,
  Palette,
} from "lucide-react"

interface CreateImageFormProps {
  onSuccess?: (id: string) => void
  onCancel?: () => void
}

const transformationIcons = {
  resize: Maximize2,
  grayscale: Palette,
  trim: Scissors,
  blur: Droplet,
  rotate: RotateCw,
}

const transformationColors = {
  resize: "from-indigo-500 to-purple-600",
  grayscale: "from-gray-500 to-gray-700",
  trim: "from-orange-500 to-red-500",
  blur: "from-cyan-500 to-blue-600",
  rotate: "from-pink-500 to-rose-600",
}

export function CreateImageForm({ onSuccess, onCancel }: CreateImageFormProps) {
  const [imageUrl, setImageUrl] = useState("")
  const [transformations, setTransformations] = useState<TransformationRequest[]>([])
  const [loading, setLoading] = useState(false)
  const [previewError, setPreviewError] = useState(false)

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

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      const result = await api.createImage({
        image_url: imageUrl,
        transformations,
      })
      onSuccess?.(result.id)
    } catch {
      alert("Failed to create image")
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Image URL Input */}
      <div className="space-y-3">
        <Label htmlFor="imageUrl" className="text-sm font-semibold text-gray-700 flex items-center gap-2">
          <Link className="w-4 h-4" />
          Image URL
        </Label>
        <Input
          id="imageUrl"
          type="url"
          value={imageUrl}
          onChange={(e) => {
            setImageUrl(e.target.value)
            setPreviewError(false)
          }}
          placeholder="https://example.com/image.jpg"
          required
          className="h-12 rounded-xl border-gray-200 focus:border-indigo-400 focus:ring-indigo-400"
        />

        {/* Image Preview */}
        {imageUrl && (
          <div className="relative rounded-2xl overflow-hidden bg-gray-100 aspect-video">
            {!previewError ? (
              <img
                src={imageUrl}
                alt="Preview"
                className="w-full h-full object-contain"
                onError={() => setPreviewError(true)}
              />
            ) : (
              <div className="w-full h-full flex items-center justify-center text-gray-400">
                <p className="text-sm">Unable to load preview</p>
              </div>
            )}
          </div>
        )}
      </div>

      {/* Transformations */}
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <Label className="text-sm font-semibold text-gray-700 flex items-center gap-2">
            <Wand2 className="w-4 h-4" />
            Transformations
            {transformations.length > 0 && (
              <span className="ml-2 px-2 py-0.5 rounded-full text-xs gradient-primary text-white">
                {transformations.length}
              </span>
            )}
          </Label>
          <Select onValueChange={(v) => addTransformation(v as TransformationRequest["name"])}>
            <SelectTrigger className="w-[180px] rounded-xl border-gray-200">
              <Plus className="w-4 h-4 mr-2" />
              <SelectValue placeholder="Add transform" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="resize">
                <span className="flex items-center gap-2">
                  <Maximize2 className="w-4 h-4" /> Resize
                </span>
              </SelectItem>
              <SelectItem value="grayscale">
                <span className="flex items-center gap-2">
                  <Palette className="w-4 h-4" /> Grayscale
                </span>
              </SelectItem>
              <SelectItem value="trim">
                <span className="flex items-center gap-2">
                  <Scissors className="w-4 h-4" /> Trim
                </span>
              </SelectItem>
              <SelectItem value="blur">
                <span className="flex items-center gap-2">
                  <Droplet className="w-4 h-4" /> Blur
                </span>
              </SelectItem>
              <SelectItem value="rotate">
                <span className="flex items-center gap-2">
                  <RotateCw className="w-4 h-4" /> Rotate
                </span>
              </SelectItem>
            </SelectContent>
          </Select>
        </div>

        {transformations.length === 0 ? (
          <Card className="p-8 text-center border-dashed border-2 border-gray-200 bg-gray-50/50 rounded-2xl">
            <Wand2 className="w-10 h-10 text-gray-300 mx-auto mb-3" />
            <p className="text-gray-500 text-sm">
              Add transformations to process your image
            </p>
          </Card>
        ) : (
          <div className="space-y-3">
            {transformations.map((transform, index) => {
              const Icon = transformationIcons[transform.name]
              const colorClass = transformationColors[transform.name]

              return (
                <Card key={index} className="p-4 rounded-2xl border-gray-100 shadow-sm hover:shadow-md transition-shadow">
                  <div className="flex items-start gap-4">
                    {/* Icon */}
                    <div className={`w-12 h-12 rounded-xl bg-gradient-to-br ${colorClass} flex items-center justify-center flex-shrink-0 shadow-lg`}>
                      <Icon className="w-6 h-6 text-white" />
                    </div>

                    {/* Content */}
                    <div className="flex-1 space-y-3">
                      <div className="flex items-center justify-between">
                        <span className="font-semibold text-gray-800 capitalize">{transform.name}</span>
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          onClick={() => removeTransformation(index)}
                          className="h-8 w-8 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors"
                        >
                          <Trash2 className="w-4 h-4" />
                        </Button>
                      </div>

                      {/* Config Fields */}
                      {transform.name === "resize" && (
                        <div className="grid grid-cols-2 gap-3">
                          <div>
                            <Label className="text-xs text-gray-500 mb-1 block">Width (px)</Label>
                            <Input
                              type="number"
                              value={transform.config.width}
                              onChange={(e) =>
                                updateTransformation(index, {
                                  ...transform.config,
                                  width: +e.target.value,
                                })
                              }
                              min={1}
                              className="h-10 rounded-lg"
                            />
                          </div>
                          <div>
                            <Label className="text-xs text-gray-500 mb-1 block">Height (px)</Label>
                            <Input
                              type="number"
                              value={transform.config.height}
                              onChange={(e) =>
                                updateTransformation(index, {
                                  ...transform.config,
                                  height: +e.target.value,
                                })
                              }
                              min={1}
                              className="h-10 rounded-lg"
                            />
                          </div>
                        </div>
                      )}

                      {transform.name === "trim" && (
                        <div>
                          <Label className="text-xs text-gray-500 mb-1 block">Threshold (0-255)</Label>
                          <Input
                            type="number"
                            value={transform.config.threshold}
                            onChange={(e) =>
                              updateTransformation(index, { threshold: +e.target.value })
                            }
                            min={0}
                            max={255}
                            className="h-10 rounded-lg"
                          />
                        </div>
                      )}

                      {transform.name === "blur" && (
                        <div>
                          <Label className="text-xs text-gray-500 mb-1 block">Sigma (blur intensity)</Label>
                          <Input
                            type="number"
                            step="0.1"
                            value={transform.config.sigma}
                            onChange={(e) =>
                              updateTransformation(index, { sigma: +e.target.value })
                            }
                            min={0.1}
                            className="h-10 rounded-lg"
                          />
                        </div>
                      )}

                      {transform.name === "rotate" && (
                        <div>
                          <Label className="text-xs text-gray-500 mb-1 block">Angle</Label>
                          <Select
                            value={transform.config.angle.toString()}
                            onValueChange={(v) => updateTransformation(index, { angle: +v })}
                          >
                            <SelectTrigger className="h-10 rounded-lg">
                              <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="90">90°</SelectItem>
                              <SelectItem value="180">180°</SelectItem>
                              <SelectItem value="270">270°</SelectItem>
                            </SelectContent>
                          </Select>
                        </div>
                      )}

                      {transform.name === "grayscale" && (
                        <p className="text-xs text-gray-400">No configuration needed</p>
                      )}
                    </div>
                  </div>
                </Card>
              )
            })}
          </div>
        )}
      </div>

      {/* Actions */}
      <div className="flex gap-3 pt-4">
        <Button
          type="submit"
          disabled={loading || !imageUrl || transformations.length === 0}
          className="flex-1 h-12 gradient-primary text-white rounded-xl font-semibold shadow-lg hover:shadow-xl transition-all disabled:opacity-50"
        >
          {loading ? (
            <>
              <Loader2 className="w-4 h-4 mr-2 animate-spin" />
              Creating...
            </>
          ) : (
            <>
              <Wand2 className="w-4 h-4 mr-2" />
              Create Image
            </>
          )}
        </Button>
        {onCancel && (
          <Button
            type="button"
            variant="outline"
            onClick={onCancel}
            className="h-12 px-6 rounded-xl border-gray-200"
          >
            Cancel
          </Button>
        )}
      </div>
    </form>
  )
}
