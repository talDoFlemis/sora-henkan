import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { api, type Image } from '@/lib/api';
import { ImageComparison } from '@/components/image-comparison';
import { CreateImageForm } from '@/components/create-image-form';
import { ArrowLeft, Trash2, Edit } from 'lucide-react';
import { env } from '@/utils/constants';

export function ImageDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [image, setImage] = useState<Image | null>(null);
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState(false);

  const loadImage = async () => {
    if (!id) return;
    try {
      const data = await api.getImage(id);
      setImage(data);
    } catch (error) {
      console.error('Failed to load image', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadImage();
  }, [id]);

  useEffect(() => {
    if (!id) return;
    const eventSource = api.streamImage(id, (updatedImage) => {
      setImage(updatedImage);
    });
    return () => eventSource.close();
  }, [id]);

  const handleDelete = async () => {
    if (!id || !confirm('Are you sure you want to delete this image?')) return;
    try {
      await api.deleteImage(id);
      navigate('/gallery');
    } catch (error) {
      alert('Failed to delete image');
    }
  };

  const handleUpdate = async (formData: any) => {
    if (!id) return;
    try {
      await api.updateImage({ id, transformations: formData.transformations });
      setEditing(false);
      loadImage();
    } catch (error) {
      alert('Failed to update image');
    }
  };

  const getImageUrl = (key: string) => `${env.AWS_BUCKET_ENDPOINT}/${key}`;

  if (loading) {
    return <div className="container mx-auto px-4 py-8">Loading...</div>;
  }

  if (!image) {
    return <div className="container mx-auto px-4 py-8">Image not found</div>;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-6">
        <Button variant="ghost" onClick={() => navigate('/gallery')}>
          <ArrowLeft className="w-4 h-4 mr-2" />
          Back to Gallery
        </Button>
      </div>

      <div className="grid lg:grid-cols-2 gap-8">
        <div className="space-y-4">
          <Card className="p-6">
            <h1 className="text-2xl font-bold mb-4">Image Details</h1>
            <dl className="space-y-2 text-sm">
              <div className="flex justify-between">
                <dt className="text-gray-600">ID:</dt>
                <dd className="font-mono">{image.id}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-gray-600">Status:</dt>
                <dd className="capitalize">{image.status}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-gray-600">MIME Type:</dt>
                <dd>{image.mime_type}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-gray-600">Transformations:</dt>
                <dd>{image.transformations.length}</dd>
              </div>
            </dl>
          </Card>

          <Card className="p-6">
            <h2 className="text-xl font-semibold mb-4">Transformations</h2>
            <ul className="space-y-2">
              {image.transformations.map((t, i) => (
                <li key={i} className="text-sm">
                  <span className="font-medium capitalize">{t.name}</span>
                  {t.name !== 'grayscale' && (
                    <span className="text-gray-600 ml-2">
                      {JSON.stringify(t.config)}
                    </span>
                  )}
                </li>
              ))}
            </ul>
          </Card>

          <div className="flex gap-3">
            <Button variant="outline" onClick={() => setEditing(!editing)}>
              <Edit className="w-4 h-4 mr-2" />
              {editing ? 'Cancel Edit' : 'Edit'}
            </Button>
            <Button variant="destructive" onClick={handleDelete}>
              <Trash2 className="w-4 h-4 mr-2" />
              Delete
            </Button>
          </div>
        </div>

        <div className="space-y-4">
          {editing ? (
            <Card className="p-6">
              <h2 className="text-xl font-semibold mb-4">Update Transformations</h2>
              <CreateImageForm
                onSuccess={() => {}}
                onCancel={() => setEditing(false)}
              />
            </Card>
          ) : (
            <ImageComparison
              beforeSrc={getImageUrl(image.object_storage_image_key)}
              afterSrc={getImageUrl(image.transformed_image_key || image.object_storage_image_key)}
              alt={image.id}
            />
          )}
        </div>
      </div>
    </div>
  );
}
