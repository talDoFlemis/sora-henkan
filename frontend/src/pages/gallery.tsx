import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { api, type Image } from '@/lib/api';
import { Plus } from 'lucide-react';
import { CreateImageForm } from '@/components/create-image-form';
import { env } from '@/utils/constants';

export function GalleryPage() {
  const navigate = useNavigate();
  const [images, setImages] = useState<Image[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);

  const loadImages = async () => {
    try {
      const response = await api.listImages(1, 50);
      setImages(response.data);
    } catch (error) {
      console.error('Failed to load images', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadImages();
  }, []);

  const handleCreateSuccess = (id: string) => {
    setShowForm(false);
    loadImages();
    navigate(`/images/${id}`);
  };

  const getImageUrl = (key: string) => `${env.AWS_BUCKET_ENDPOINT}/${key}`;

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-4xl font-bold">Image Gallery</h1>
        <Button onClick={() => setShowForm(!showForm)}>
          <Plus className="w-4 h-4 mr-2" />
          Create New Image
        </Button>
      </div>

      {showForm && (
        <Card className="p-6 mb-8">
          <h2 className="text-2xl font-semibold mb-4">Create New Image</h2>
          <CreateImageForm onSuccess={handleCreateSuccess} onCancel={() => setShowForm(false)} />
        </Card>
      )}

      {loading ? (
        <div className="text-center py-12">Loading...</div>
      ) : images.length === 0 ? (
        <div className="text-center py-12 text-gray-500">
          No images yet. Create your first one!
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {images.map((image) => (
            <Card
              key={image.id}
              className="cursor-pointer hover:shadow-lg transition-shadow overflow-hidden"
              onClick={() => navigate(`/images/${image.id}`)}
            >
              <div className="aspect-video bg-gray-100">
                <img
                  src={getImageUrl(image.transformed_image_key || image.object_storage_image_key)}
                  alt={image.id}
                  className="w-full h-full object-cover"
                />
              </div>
              <div className="p-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium truncate">{image.id.slice(0, 8)}...</span>
                  <span className={`text-xs px-2 py-1 rounded ${
                    image.status === 'completed' ? 'bg-green-100 text-green-800' :
                    image.status === 'processing' ? 'bg-yellow-100 text-yellow-800' :
                    'bg-gray-100 text-gray-800'
                  }`}>
                    {image.status}
                  </span>
                </div>
                <p className="text-xs text-gray-500 mt-2">
                  {image.transformations.length} transformation(s)
                </p>
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
