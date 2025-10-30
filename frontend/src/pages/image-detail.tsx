import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { api, type Image, type TransformationRequest } from '@/lib/api';
import { ImageComparison } from '@/components/image-comparison';
import { ArrowLeft, Trash2, Edit, Save, X, Plus, Loader2, ExternalLink, Copy, Check, ImageIcon } from 'lucide-react';
import { env } from '@/utils/constants';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import { Trash2 as TrashIcon } from 'lucide-react';

export function ImageDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [image, setImage] = useState<Image | null>(null);
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState(false);
  const [transformations, setTransformations] = useState<TransformationRequest[]>([]);
  const [updating, setUpdating] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [copiedUrl, setCopiedUrl] = useState(false);

  const loadImage = async () => {
    if (!id) return;
    setLoading(true);
    try {
      const data = await api.getImage(id);
      setImage(data);
      setTransformations(data.transformations);
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
    if (!id || !confirm('Are you sure you want to delete this image? This action cannot be undone.')) return;
    setDeleting(true);
    try {
      await api.deleteImage(id);
      navigate('/gallery');
    } catch (error) {
      alert('Failed to delete image');
      setDeleting(false);
    }
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!id) return;
    setUpdating(true);
    try {
      await api.updateImage({ id, transformations });
      setEditing(false);
      await loadImage();
    } catch (error) {
      alert('Failed to update image');
    } finally {
      setUpdating(false);
    }
  };

  const addTransformation = (type: TransformationRequest['name']) => {
    const newTransform: TransformationRequest = 
      type === 'resize' ? { name: 'resize', config: { width: 800, height: 600 } } :
      type === 'grayscale' ? { name: 'grayscale', config: {} } :
      type === 'trim' ? { name: 'trim', config: { threshold: 10 } } :
      type === 'blur' ? { name: 'blur', config: { sigma: 1.5 } } :
      { name: 'rotate', config: { angle: 90 } };
    
    setTransformations([...transformations, newTransform]);
  };

  const updateTransformation = (index: number, config: any) => {
    const updated = [...transformations];
    updated[index] = { ...updated[index], config };
    setTransformations(updated);
  };

  const removeTransformation = (index: number) => {
    setTransformations(transformations.filter((_, i) => i !== index));
  };

  const getImageUrl = (key: string) => `${env.AWS_BUCKET_ENDPOINT}/${key}`;

  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
      setCopiedUrl(true);
      setTimeout(() => setCopiedUrl(false), 2000);
    } catch (error) {
      console.error('Failed to copy', error);
    }
  };

  if (loading) {
    return (
      <div className=" flex items-center justify-center">
        <div className="text-center">
          <Loader2 className="w-12 h-12 animate-spin text-indigo-600 mx-auto mb-4" />
          <p className="text-gray-600">Loading image details...</p>
        </div>
      </div>
    );
  }

  if (!image) {
    return (
      <div className=" flex items-center justify-center">
        <Card className="p-8 text-center max-w-md">
          <ImageIcon className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h2 className="text-2xl font-bold mb-2">Image Not Found</h2>
          <p className="text-gray-600 mb-6">The image you're looking for doesn't exist or has been deleted.</p>
          <Button onClick={() => navigate('/gallery')} className="">
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back to Gallery
          </Button>
        </Card>
      </div>
    );
  }

  const statusColors = {
    completed: 'bg-green-500 text-white',
    processing: 'bg-yellow-500 text-white',
    pending: 'bg-blue-500 text-white',
    failed: 'bg-red-500 text-white',
  };

  return (
    <div className="">
      <div className="container mx-auto px-4 py-8">
        {/* Header */}
        <div className="mb-8">
          <Button 
           
            onClick={() => navigate('/gallery')}
            className="mb-4 hover:bg-white/60"
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back to Gallery
          </Button>

          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <div>
              <div className="flex items-center gap-3 mb-2">
                <h1 className="text-4xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
                  Image Details
                </h1>
                <Badge className={statusColors[image.status as keyof typeof statusColors] || 'bg-gray-500'}>
                  {image.status}
                </Badge>
              </div>
              <p className="text-gray-600 font-mono text-sm">{image.id}</p>
            </div>

            <div className="flex gap-3">
              <Button 
                variant="outline"
                onClick={() => setEditing(!editing)}
                disabled={editing && updating}
                className="border-indigo-200 text-indigo-700 hover:bg-indigo-50"
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
                className="bg-red-600 hover:bg-red-700 text-white"
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
        </div>

        {/* Main Content */}
        <div className="grid lg:grid-cols-3 gap-8">
          {/* Left Column - Image Comparison */}
          <div className="lg:col-span-2 space-y-6">
            {editing ? (
              <Card className="p-6 bg-white/80 backdrop-blur-sm">
                <h2 className="text-2xl font-bold mb-6">Edit Transformations</h2>
                <form onSubmit={handleUpdate} className="space-y-6">
                  <div className="space-y-4">
                    <div className="flex items-center justify-between">
                      <Label className="text-lg font-semibold">Transformations</Label>
                      <Select onValueChange={(v) => addTransformation(v as any)}>
                        <SelectTrigger className="w-[220px]">
                          <Plus className="w-4 h-4 mr-2" />
                          <SelectValue placeholder="Add transformation" />
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

                    {transformations.length === 0 && (
                      <div className="text-center py-8 text-gray-500">
                        <ImageIcon className="w-12 h-12 mx-auto mb-3 opacity-50" />
                        <p>No transformations added yet</p>
                        <p className="text-sm">Add at least one transformation above</p>
                      </div>
                    )}

                    {transformations.map((transform, index) => (
                      <Card key={index} className="p-4 bg-gradient-to-br from-gray-50 to-gray-100 border-gray-200">
                        <div className="flex items-start justify-between gap-4">
                          <div className="flex-1 space-y-3">
                            <div className="flex items-center gap-2">
                              <Badge variant="outline" className="capitalize font-semibold mb-2">
                                {transform.name}
                              </Badge>
                            </div>
                            {transform.name === 'resize' && (
                              <div className="grid grid-cols-2 gap-3">
                                <div>
                                  <Label className="text-sm">Width (px)</Label>
                                  <Input
                                    type="number"
                                    value={transform.config.width}
                                    onChange={(e) => updateTransformation(index, { ...transform.config, width: +e.target.value })}
                                    min={1}
                                    className="mt-1"
                                  />
                                </div>
                                <div>
                                  <Label className="text-sm">Height (px)</Label>
                                  <Input
                                    type="number"
                                    value={transform.config.height}
                                    onChange={(e) => updateTransformation(index, { ...transform.config, height: +e.target.value })}
                                    min={1}
                                    className="mt-1"
                                  />
                                </div>
                              </div>
                            )}
                            {transform.name === 'trim' && (
                              <div>
                                <Label className="text-sm">Threshold (0-255)</Label>
                                <Input
                                  type="number"
                                  value={transform.config.threshold}
                                  onChange={(e) => updateTransformation(index, { threshold: +e.target.value })}
                                  min={0}
                                  max={255}
                                  className="mt-1"
                                />
                              </div>
                            )}
                            {transform.name === 'blur' && (
                              <div>
                                <Label className="text-sm">Sigma (Blur Intensity)</Label>
                                <Input
                                  type="number"
                                  step="0.1"
                                  value={transform.config.sigma}
                                  onChange={(e) => updateTransformation(index, { sigma: +e.target.value })}
                                  min={0.1}
                                  className="mt-1"
                                />
                              </div>
                            )}
                            {transform.name === 'rotate' && (
                              <div>
                                <Label className="text-sm">Rotation Angle</Label>
                                <Select
                                  value={transform.config.angle.toString()}
                                  onValueChange={(v) => updateTransformation(index, { angle: +v })}
                                >
                                  <SelectTrigger className="mt-1">
                                    <SelectValue />
                                  </SelectTrigger>
                                  <SelectContent>
                                    <SelectItem value="90">90° (Quarter turn)</SelectItem>
                                    <SelectItem value="180">180° (Half turn)</SelectItem>
                                    <SelectItem value="270">270° (Three-quarter turn)</SelectItem>
                                  </SelectContent>
                                </Select>
                              </div>
                            )}
                            {transform.name === 'grayscale' && (
                              <p className="text-sm text-gray-600">No configuration required</p>
                            )}
                          </div>
                          <Button 
                            type="button" 
                            variant="ghost" 
                            size="icon" 
                            onClick={() => removeTransformation(index)}
                            className="text-red-600 hover:text-red-700 hover:bg-red-50"
                          >
                            <TrashIcon className="w-4 h-4" />
                          </Button>
                        </div>
                      </Card>
                    ))}
                  </div>

                  <div className="flex gap-3 pt-4 border-t">
                    <Button 
                      type="submit" 
                      disabled={updating || transformations.length === 0}
                      className="bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white"
                    >
                      {updating ? (
                        <>
                          <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                          Updating...
                        </>
                      ) : (
                        <>
                          <Save className="w-4 h-4 mr-2" />
                          Save Changes
                        </>
                      )}
                    </Button>
                    <Button 
                      type="button" 
                      variant="outline" 
                      onClick={() => {
                        setEditing(false);
                        setTransformations(image.transformations);
                      }}
                      disabled={updating}
                    >
                      Cancel
                    </Button>
                  </div>
                </form>
              </Card>
            ) : (
              <Card className="p-6 bg-white/80 backdrop-blur-sm overflow-hidden">
                <h2 className="text-xl font-bold mb-4">Image Comparison</h2>
                <ImageComparison
                  beforeSrc={getImageUrl(image.object_storage_image_key)}
                  afterSrc={getImageUrl(image.transformed_image_key || image.object_storage_image_key)}
                  alt={image.id}
                />
              </Card>
            )}
          </div>

          {/* Right Column - Details */}
          <div className="space-y-6">
            {/* Image Information */}
            <Card className="p-6 bg-white/80 backdrop-blur-sm">
              <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
                <ImageIcon className="w-5 h-5 text-indigo-600" />
                Information
              </h2>
              <dl className="space-y-3">
                <div className="flex flex-col gap-1">
                  <dt className="text-sm text-gray-600 font-medium">Image ID</dt>
                  <dd className="font-mono text-sm bg-gray-100 p-2 rounded flex items-center justify-between gap-2">
                    <span className="truncate">{image.id}</span>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-6 w-6 flex-shrink-0"
                      onClick={() => copyToClipboard(image.id)}
                    >
                      {copiedUrl ? (
                        <Check className="w-3 h-3 text-green-600" />
                      ) : (
                        <Copy className="w-3 h-3" />
                      )}
                    </Button>
                  </dd>
                </div>
                <div className="flex justify-between py-2 border-b">
                  <dt className="text-sm text-gray-600">Status</dt>
                  <dd>
                    <Badge className={statusColors[image.status as keyof typeof statusColors] || 'bg-gray-500'}>
                      {image.status}
                    </Badge>
                  </dd>
                </div>
                <div className="flex justify-between py-2 border-b">
                  <dt className="text-sm text-gray-600">MIME Type</dt>
                  <dd className="text-sm font-medium">{image.mime_type}</dd>
                </div>
                <div className="flex justify-between py-2 border-b">
                  <dt className="text-sm text-gray-600">Checksum</dt>
                  <dd className="text-sm font-mono truncate max-w-[150px]" title={image.checksum}>
                    {image.checksum.slice(0, 12)}...
                  </dd>
                </div>
                <div className="flex justify-between py-2 border-b">
                  <dt className="text-sm text-gray-600">Created</dt>
                  <dd className="text-sm">
                    {new Date(image.created_at).toLocaleDateString('en-US', {
                      month: 'short',
                      day: 'numeric',
                      year: 'numeric',
                      hour: '2-digit',
                      minute: '2-digit',
                    })}
                  </dd>
                </div>
                <div className="flex justify-between py-2">
                  <dt className="text-sm text-gray-600">Updated</dt>
                  <dd className="text-sm">
                    {new Date(image.updated_at).toLocaleDateString('en-US', {
                      month: 'short',
                      day: 'numeric',
                      year: 'numeric',
                      hour: '2-digit',
                      minute: '2-digit',
                    })}
                  </dd>
                </div>
              </dl>
            </Card>

            {/* URLs */}
            <Card className="p-6 bg-white/80 backdrop-blur-sm">
              <h2 className="text-xl font-bold mb-4">Image URLs</h2>
              <div className="space-y-4">
                <div>
                  <Label className="text-sm text-gray-600 mb-2 block">Original Image</Label>
                  <div className="flex gap-2">
                    <Input
                      value={image.original_image_url}
                      readOnly
                      className="text-sm font-mono"
                    />
                    <Button
                      variant="outline"
                      size="icon"
                      onClick={() => window.open(image.original_image_url, '_blank')}
                    >
                      <ExternalLink className="w-4 h-4" />
                    </Button>
                  </div>
                </div>
                <div>
                  <Label className="text-sm text-gray-600 mb-2 block">Transformed Image</Label>
                  <div className="flex gap-2">
                    <Input
                      value={getImageUrl(image.transformed_image_key || image.object_storage_image_key)}
                      readOnly
                      className="text-sm font-mono"
                    />
                    <Button
                      variant="outline"
                      size="icon"
                      onClick={() => window.open(getImageUrl(image.transformed_image_key || image.object_storage_image_key), '_blank')}
                    >
                      <ExternalLink className="w-4 h-4" />
                    </Button>
                  </div>
                </div>
              </div>
            </Card>

            {/* Transformations */}
            <Card className="p-6 bg-white/80 backdrop-blur-sm">
              <h2 className="text-xl font-bold mb-4">
                Applied Transformations ({image.transformations.length})
              </h2>
              {image.transformations.length === 0 ? (
                <p className="text-gray-500 text-sm text-center py-4">No transformations applied</p>
              ) : (
                <ul className="space-y-3">
                  {image.transformations.map((t, i) => (
                    <li key={i} className="bg-gradient-to-r from-indigo-50 to-purple-50 p-3 rounded-lg border border-indigo-100">
                      <div className="flex items-start justify-between gap-2">
                        <div className="flex-1">
                          <Badge variant="outline" className="capitalize font-semibold mb-2">
                            {i + 1}. {t.name}
                          </Badge>
                          {t.name !== 'grayscale' && (
                            <div className="mt-2 text-xs text-gray-700 font-mono bg-white/60 p-2 rounded">
                              {JSON.stringify(t.config, null, 2)}
                            </div>
                          )}
                        </div>
                      </div>
                    </li>
                  ))}
                </ul>
              )}
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}
