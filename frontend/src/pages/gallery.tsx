import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { api, type Image } from "@/lib/api"
import {
  Plus,
  Search,
  Filter,
  Grid3x3,
  LayoutGrid,
  X,
  Loader2,
  Home,
  Sparkles,
  ImageIcon,
} from "lucide-react"
import { CreateImageForm } from "@/components/create-image-form"
import { env } from "@/utils/constants"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Badge } from "@/components/ui/badge"

export function GalleryPage() {
  const navigate = useNavigate()
  const [images, setImages] = useState<Image[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [searchQuery, setSearchQuery] = useState("")
  const [statusFilter, setStatusFilter] = useState<string>("all")
  const [gridSize, setGridSize] = useState<"small" | "medium" | "large">(
    "medium",
  )
  const [page, setPage] = useState(1)
  const [totalCount, setTotalCount] = useState(0)
  const limit = 12

  const loadImages = async () => {
    setLoading(true)
    try {
      const response = await api.listImages(page, limit)
      setImages(response.data)
      setTotalCount(response.count)
    } catch (error) {
      console.error("Failed to load images", error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadImages()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page])

  const handleCreateSuccess = (id: string) => {
    setShowForm(false)
    loadImages()
    navigate(`/images/${id}`)
  }

  const getImageUrl = (key: string) => `${env.AWS_BUCKET_ENDPOINT}/${key}`

  const filteredImages = images.filter((image) => {
    const matchesSearch = image.id
      .toLowerCase()
      .includes(searchQuery.toLowerCase())
    const matchesStatus =
      statusFilter === "all" || image.status === statusFilter
    return matchesSearch && matchesStatus
  })

  const gridCols = {
    small: "grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5",
    medium: "grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4",
    large: "grid-cols-1 md:grid-cols-2 lg:grid-cols-3",
  }

  const totalPages = Math.ceil(totalCount / limit)

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
                onClick={() => setShowForm(true)}
                className="gradient-primary text-white rounded-xl shadow-lg hover:shadow-xl transition-all hover:-translate-y-0.5"
              >
                <Plus className="w-4 h-4 mr-2" />
                New Image
              </Button>
            </div>
          </div>
        </div>
      </header>

      <main className="w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Page Title */}
        <div className="mb-8 animate-fade-in-up">
          <h1 className="text-4xl md:text-5xl font-bold text-gray-800 mb-2">
            Image <span className="gradient-text">Gallery</span>
          </h1>
          <p className="text-gray-600">
            {totalCount} {totalCount === 1 ? "image" : "images"} • Page {page}{" "}
            of {totalPages || 1}
          </p>
        </div>

        {/* Filters */}
        <Card
          className="p-5 mb-8 glass-dark rounded-2xl shadow-lg border-0 animate-fade-in-up"
          style={{ animationDelay: "0.1s" }}
        >
          <div className="flex flex-col lg:flex-row gap-4">
            {/* Search */}
            <div className="flex-1 relative">
              <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
              <Input
                placeholder="Search by image ID..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-12 h-12 rounded-xl border-gray-200 focus:border-indigo-400 focus:ring-indigo-400"
              />
            </div>

            {/* Status Filter */}
            <div className="w-full lg:w-52">
              <Select value={statusFilter} onValueChange={setStatusFilter}>
                <SelectTrigger className="h-12 rounded-xl border-gray-200">
                  <Filter className="w-4 h-4 mr-2 text-gray-400" />
                  <SelectValue placeholder="Filter by status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="pending">Pending</SelectItem>
                  <SelectItem value="processing">Processing</SelectItem>
                  <SelectItem value="processed">Processed</SelectItem>
                  <SelectItem value="failed">Failed</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* Grid Size */}
            <div className="flex gap-2">
              <Button
                variant={gridSize === "small" ? "default" : "outline"}
                size="icon"
                onClick={() => setGridSize("small")}
                className={`h-12 w-12 rounded-xl ${gridSize === "small" ? "gradient-primary text-white" : ""}`}
              >
                <Grid3x3 className="w-4 h-4" />
              </Button>
              <Button
                variant={gridSize === "medium" ? "default" : "outline"}
                size="icon"
                onClick={() => setGridSize("medium")}
                className={`h-12 w-12 rounded-xl ${gridSize === "medium" ? "gradient-primary text-white" : ""}`}
              >
                <LayoutGrid className="w-4 h-4" />
              </Button>
              <Button
                variant={gridSize === "large" ? "default" : "outline"}
                size="icon"
                onClick={() => setGridSize("large")}
                className={`h-12 w-12 rounded-xl ${gridSize === "large" ? "gradient-primary text-white" : ""}`}
              >
                <LayoutGrid className="w-5 h-5" />
              </Button>
            </div>
          </div>

          {/* Active Filters */}
          {(searchQuery || statusFilter !== "all") && (
            <div className="flex flex-wrap gap-2 mt-4 pt-4 border-t border-gray-200">
              <span className="text-sm text-gray-500">Active filters:</span>
              {searchQuery && (
                <Badge
                  variant="secondary"
                  className="gap-1 rounded-lg px-3 py-1"
                >
                  Search: {searchQuery}
                  <X
                    className="w-3 h-3 cursor-pointer hover:text-red-500 transition-colors"
                    onClick={() => setSearchQuery("")}
                  />
                </Badge>
              )}
              {statusFilter !== "all" && (
                <Badge
                  variant="secondary"
                  className="gap-1 rounded-lg px-3 py-1"
                >
                  Status: {statusFilter}
                  <X
                    className="w-3 h-3 cursor-pointer hover:text-red-500 transition-colors"
                    onClick={() => setStatusFilter("all")}
                  />
                </Badge>
              )}
            </div>
          )}
        </Card>

        {/* Content */}
        {loading ? (
          <div className="flex flex-col items-center justify-center py-24">
            <div className="w-16 h-16 rounded-2xl gradient-primary flex items-center justify-center mb-6 animate-pulse">
              <Loader2 className="w-8 h-8 animate-spin text-white" />
            </div>
            <p className="text-gray-600 font-medium">Loading your gallery...</p>
          </div>
        ) : filteredImages.length === 0 ? (
          <Card className="p-16 text-center glass-dark rounded-3xl border-0 shadow-lg animate-fade-in-up">
            <div className="max-w-md mx-auto">
              <div className="w-24 h-24 rounded-3xl gradient-primary flex items-center justify-center mx-auto mb-6 shadow-lg">
                <ImageIcon className="w-12 h-12 text-white" />
              </div>
              <h3 className="text-2xl font-bold text-gray-800 mb-3">
                No images found
              </h3>
              <p className="text-gray-600 mb-8">
                {searchQuery || statusFilter !== "all"
                  ? "Try adjusting your filters or search query"
                  : "Get started by creating your first image transformation"}
              </p>
              {!searchQuery && statusFilter === "all" && (
                <Button
                  onClick={() => setShowForm(true)}
                  className="gradient-primary text-white rounded-xl px-8 py-6 text-lg shadow-lg hover:shadow-xl transition-all"
                >
                  <Plus className="w-5 h-5 mr-2" />
                  Create Your First Image
                </Button>
              )}
            </div>
          </Card>
        ) : (
          <>
            {/* Image Grid */}
            <div className={`grid ${gridCols[gridSize]} gap-6 mb-8`}>
              {filteredImages.map((image, index) => (
                <Card
                  key={image.id}
                  className="group cursor-pointer overflow-hidden glass-dark rounded-2xl border-0 shadow-lg hover-lift animate-fade-in-up"
                  style={{ animationDelay: `${index * 0.05}s` }}
                  onClick={() => navigate(`/images/${image.id}`)}
                >
                  {/* Image */}
                  <div className="aspect-square bg-gradient-to-br from-gray-100 to-gray-200 overflow-hidden relative">
                    <img
                      src={getImageUrl(
                        image.transformed_image_key ||
                          image.object_storage_image_key,
                      )}
                      alt={image.id}
                      className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
                      loading="lazy"
                    />
                    {/* Overlay */}
                    <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/0 to-black/0 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                      <div className="absolute bottom-0 left-0 right-0 p-4 text-white transform translate-y-4 group-hover:translate-y-0 transition-transform duration-300">
                        <p className="text-sm font-medium">
                          Click to view details
                        </p>
                      </div>
                    </div>
                    {/* Status Badge */}
                    <div className="absolute top-3 right-3">
                      <Badge
                        className={`${getStatusColor(image.status)} text-white shadow-lg border-0 rounded-lg px-3`}
                      >
                        {image.status}
                      </Badge>
                    </div>
                  </div>

                  {/* Info */}
                  <div className="p-4">
                    <div className="flex items-start justify-between gap-2 mb-2">
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-mono font-bold text-gray-800 truncate">
                          {image.id.slice(0, 8)}...
                        </p>
                        <p className="text-xs text-gray-500 mt-1">
                          {new Date(image.created_at).toLocaleDateString(
                            "en-US",
                            {
                              month: "short",
                              day: "numeric",
                              year: "numeric",
                            },
                          )}
                        </p>
                      </div>
                    </div>

                    {/* Transformations */}
                    {image.transformations.length > 0 && (
                      <div className="flex flex-wrap gap-1.5 mt-3">
                        {image.transformations.slice(0, 3).map((t, i) => (
                          <Badge
                            key={i}
                            variant="outline"
                            className="text-xs capitalize rounded-md border-indigo-200 text-indigo-600"
                          >
                            {t.name}
                          </Badge>
                        ))}
                        {image.transformations.length > 3 && (
                          <Badge
                            variant="outline"
                            className="text-xs rounded-md"
                          >
                            +{image.transformations.length - 3}
                          </Badge>
                        )}
                      </div>
                    )}
                  </div>
                </Card>
              ))}
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <Card className="p-4 glass-dark rounded-2xl border-0 shadow-lg">
                <div className="flex items-center justify-between">
                  <Button
                    variant="outline"
                    onClick={() => setPage((p) => Math.max(1, p - 1))}
                    disabled={page === 1}
                    className="rounded-xl"
                  >
                    Previous
                  </Button>

                  <div className="flex items-center gap-2">
                    {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                      const pageNum = i + 1
                      return (
                        <Button
                          key={pageNum}
                          variant={page === pageNum ? "default" : "outline"}
                          size="icon"
                          onClick={() => setPage(pageNum)}
                          className={`rounded-xl ${page === pageNum ? "gradient-primary text-white" : ""}`}
                        >
                          {pageNum}
                        </Button>
                      )
                    })}
                    {totalPages > 5 && (
                      <span className="text-gray-400 px-2">...</span>
                    )}
                  </div>

                  <Button
                    variant="outline"
                    onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                    disabled={page === totalPages}
                    className="rounded-xl"
                  >
                    Next
                  </Button>
                </div>
              </Card>
            )}
          </>
        )}
      </main>

      {/* Create Image Dialog */}
      <Dialog open={showForm} onOpenChange={setShowForm}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto rounded-3xl">
          <DialogHeader>
            <DialogTitle className="text-2xl font-bold flex items-center gap-3">
              <div className="w-10 h-10 rounded-xl gradient-primary flex items-center justify-center">
                <Plus className="w-5 h-5 text-white" />
              </div>
              Create New Image
            </DialogTitle>
            <DialogDescription className="text-gray-600">
              Upload an image URL and apply transformations to create a new
              processed image.
            </DialogDescription>
          </DialogHeader>
          <CreateImageForm
            onSuccess={handleCreateSuccess}
            onCancel={() => setShowForm(false)}
          />
        </DialogContent>
      </Dialog>
    </div>
  )
}
