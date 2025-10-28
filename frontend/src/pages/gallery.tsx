import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { api, type Image } from '@/lib/api';
import { Plus, Search, Filter, Grid3x3, LayoutGrid, X, Loader2 } from 'lucide-react';
import { CreateImageForm } from '@/components/create-image-form';
import { env } from '@/utils/constants';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';

export function GalleryPage() {
  const navigate = useNavigate();
  const [images, setImages] = useState<Image[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [gridSize, setGridSize] = useState<'small' | 'medium' | 'large'>('medium');
  const [page, setPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const limit = 12;

  const loadImages = async () => {
    setLoading(true);
    try {
      const response = await api.listImages(page, limit);
      setImages(response.data);
      setTotalCount(response.count);
    } catch (error) {
      console.error('Failed to load images', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadImages();
  }, [page]);

  const handleCreateSuccess = (id: string) => {
    setShowForm(false);
    loadImages();
    navigate(`/images/${id}`);
  };

  const getImageUrl = (key: string) => `${env.AWS_BUCKET_ENDPOINT}/${key}`;

  const filteredImages = images.filter((image) => {
    const matchesSearch = image.id.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesStatus = statusFilter === 'all' || image.status === statusFilter;
    return matchesSearch && matchesStatus;
  });

  const gridCols = {
    small: 'grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5',
    medium: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4',
    large: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3',
  };

  const totalPages = Math.ceil(totalCount / limit);

  return (
    <div className="">
      <div className=" px-4 py-8">
        {/* Header */}
        <div className="mb-8">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6">
            <div>
              <h1 className="text-4xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent mb-2">
                Image Gallery
              </h1>
              <p className="text-gray-600">
                {totalCount} {totalCount === 1 ? 'image' : 'images'} • Page {page} of {totalPages || 1}
              </p>
            </div>
            <Button 
              onClick={() => setShowForm(true)}
              className="bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white shadow-lg hover:shadow-xl transition-all"
            >
              <Plus className="w-4 h-4 mr-2" />
              Create New Image
            </Button>
          </div>

          {/* Filters & Search Bar */}
          <Card className="p-4 bg-white/80 backdrop-blur-sm border-gray-200">
            <div className="flex flex-col lg:flex-row gap-4">
              {/* Search */}
              <div className="flex-1 relative">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                <Input
                  placeholder="Search by image ID..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10"
                />
              </div>

              {/* Status Filter */}
              <div className="w-full lg:w-48">
                <Select value={statusFilter} onValueChange={setStatusFilter}>
                  <SelectTrigger>
                    <Filter className="w-4 h-4 mr-2" />
                    <SelectValue placeholder="Filter by status" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Status</SelectItem>
                    <SelectItem value="pending">Pending</SelectItem>
                    <SelectItem value="processing">Processing</SelectItem>
                    <SelectItem value="completed">Completed</SelectItem>
                    <SelectItem value="failed">Failed</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Grid Size Toggle */}
              <div className="flex gap-2">
                <Button
                  variant={gridSize === 'small' ? 'default' : 'outline'}
                  size="icon"
                  onClick={() => setGridSize('small')}
                  title="Small grid"
                >
                  <Grid3x3 className="w-4 h-4" />
                </Button>
                <Button
                  variant={gridSize === 'medium' ? 'default' : 'outline'}
                  size="icon"
                  onClick={() => setGridSize('medium')}
                  title="Medium grid"
                >
                  <LayoutGrid className="w-4 h-4" />
                </Button>
                <Button
                  variant={gridSize === 'large' ? 'default' : 'outline'}
                  size="icon"
                  onClick={() => setGridSize('large')}
                  title="Large grid"
                >
                  <LayoutGrid className="w-5 h-5" />
                </Button>
              </div>
            </div>

            {/* Active Filters */}
            {(searchQuery || statusFilter !== 'all') && (
              <div className="flex flex-wrap gap-2 mt-4 pt-4 border-t">
                <span className="text-sm text-gray-600">Active filters:</span>
                {searchQuery && (
                  <Badge variant="secondary" className="gap-1">
                    Search: {searchQuery}
                    <X
                      className="w-3 h-3 cursor-pointer hover:text-red-600"
                      onClick={() => setSearchQuery('')}
                    />
                  </Badge>
                )}
                {statusFilter !== 'all' && (
                  <Badge variant="secondary" className="gap-1">
                    Status: {statusFilter}
                    <X
                      className="w-3 h-3 cursor-pointer hover:text-red-600"
                      onClick={() => setStatusFilter('all')}
                    />
                  </Badge>
                )}
              </div>
            )}
          </Card>
        </div>

        {/* Content */}
        {loading ? (
          <div className="flex flex-col items-center justify-center py-20">
            <Loader2 className="w-12 h-12 animate-spin text-indigo-600 mb-4" />
            <p className="text-gray-600">Loading images...</p>
          </div>
        ) : filteredImages.length === 0 ? (
          <Card className="p-12 text-center bg-white/60 backdrop-blur-sm">
            <div className="max-w-md mx-auto">
              <div className="w-20 h-20 bg-gradient-to-br from-indigo-100 to-purple-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <Search className="w-10 h-10 text-indigo-600" />
              </div>
              <h3 className="text-xl font-semibold mb-2">No images found</h3>
              <p className="text-gray-600 mb-6">
                {searchQuery || statusFilter !== 'all'
                  ? 'Try adjusting your filters or search query'
                  : 'Get started by creating your first image transformation'}
              </p>
              {!searchQuery && statusFilter === 'all' && (
                <Button onClick={() => setShowForm(true)} className="bg-indigo-600 hover:bg-indigo-700">
                  <Plus className="w-4 h-4 mr-2" />
                  Create Your First Image
                </Button>
              )}
            </div>
          </Card>
        ) : (
          <>
            {/* Image Grid */}
            <div className={`grid ${gridCols[gridSize]} gap-6 mb-8`}>
              {filteredImages.map((image) => (
                <Card
                  key={image.id}
                  className="group cursor-pointer hover:shadow-2xl transition-all duration-300 overflow-hidden bg-white border-gray-200 hover:border-indigo-300"
                  onClick={() => navigate(`/images/${image.id}`)}
                >
                  {/* Image */}
                  <div className="aspect-square bg-gradient-to-br from-gray-100 to-gray-200 overflow-hidden relative">
                    <img
                      src={getImageUrl(image.transformed_image_key || image.object_storage_image_key)}
                      alt={image.id}
                      className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                      loading="lazy"
                    />
                    {/* Overlay on hover */}
                    <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-black/0 to-black/0 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                      <div className="absolute bottom-0 left-0 right-0 p-4 text-white">
                        <p className="text-sm font-medium truncate">Click to view details</p>
                      </div>
                    </div>
                    {/* Status Badge */}
                    <div className="absolute top-3 right-3">
                      <Badge
                        className={`
                          ${image.status === 'completed' ? 'bg-green-500 hover:bg-green-600' : ''}
                          ${image.status === 'processing' ? 'bg-yellow-500 hover:bg-yellow-600' : ''}
                          ${image.status === 'pending' ? 'bg-blue-500 hover:bg-blue-600' : ''}
                          ${image.status === 'failed' ? 'bg-red-500 hover:bg-red-600' : ''}
                          text-white shadow-lg
                        `}
                      >
                        {image.status}
                      </Badge>
                    </div>
                  </div>

                  {/* Info */}
                  <div className="p-4">
                    <div className="flex items-start justify-between gap-2 mb-2">
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-mono font-semibold truncate text-gray-900">
                          {image.id.slice(0, 8)}...
                        </p>
                        <p className="text-xs text-gray-500 mt-1">
                          {new Date(image.created_at).toLocaleDateString('en-US', {
                            month: 'short',
                            day: 'numeric',
                            year: 'numeric',
                          })}
                        </p>
                      </div>
                    </div>
                    
                    <div className="flex items-center gap-2 text-xs text-gray-600">
                      <span className="flex items-center gap-1">
                        <div className="w-1.5 h-1.5 rounded-full bg-indigo-600"></div>
                        {image.transformations.length} transformation{image.transformations.length !== 1 ? 's' : ''}
                      </span>
                    </div>

                    {/* Transformation Tags */}
                    {image.transformations.length > 0 && (
                      <div className="flex flex-wrap gap-1 mt-3">
                        {image.transformations.slice(0, 3).map((t, i) => (
                          <Badge key={i} variant="outline" className="text-xs capitalize">
                            {t.name}
                          </Badge>
                        ))}
                        {image.transformations.length > 3 && (
                          <Badge variant="outline" className="text-xs">
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
              <Card className="p-4 bg-white/80 backdrop-blur-sm border-gray-200">
                <div className="flex items-center justify-between">
                  <Button
                    variant="outline"
                    onClick={() => setPage((p) => Math.max(1, p - 1))}
                    disabled={page === 1}
                  >
                    Previous
                  </Button>
                  
                  <div className="flex items-center gap-2">
                    {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                      const pageNum = i + 1;
                      return (
                        <Button
                          key={pageNum}
                          variant={page === pageNum ? 'default' : 'outline'}
                          size="icon"
                          onClick={() => setPage(pageNum)}
                          className={page === pageNum ? 'bg-indigo-600 hover:bg-indigo-700' : ''}
                        >
                          {pageNum}
                        </Button>
                      );
                    })}
                    {totalPages > 5 && <span className="text-gray-500">...</span>}
                  </div>

                  <Button
                    variant="outline"
                    onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                    disabled={page === totalPages}
                  >
                    Next
                  </Button>
                </div>
              </Card>
            )}
          </>
        )}
      </div>

      {/* Create Image Dialog */}
      <Dialog open={showForm} onOpenChange={setShowForm}>
        <DialogContent className="max-w-3xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle className="text-2xl font-bold">Create New Image</DialogTitle>
            <DialogDescription>
              Upload an image URL and apply transformations to create a new processed image.
            </DialogDescription>
          </DialogHeader>
          <CreateImageForm onSuccess={handleCreateSuccess} onCancel={() => setShowForm(false)} />
        </DialogContent>
      </Dialog>
    </div>
  );
}
