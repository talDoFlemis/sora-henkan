import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Card } from '@/components/ui/card';
import { api, type TransformationRequest } from '@/lib/api';
import { Plus, Trash2 } from 'lucide-react';

interface CreateImageFormProps {
  onSuccess?: (id: string) => void;
  onCancel?: () => void;
}

export function CreateImageForm({ onSuccess, onCancel }: CreateImageFormProps) {
  const [imageUrl, setImageUrl] = useState('');
  const [transformations, setTransformations] = useState<TransformationRequest[]>([]);
  const [loading, setLoading] = useState(false);

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

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      const result = await api.createImage({ image_url: imageUrl, transformations });
      onSuccess?.(result.id);
    } catch (error) {
      alert('Failed to create image');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div>
        <Label htmlFor="imageUrl">Image URL</Label>
        <Input
          id="imageUrl"
          type="url"
          value={imageUrl}
          onChange={(e) => setImageUrl(e.target.value)}
          placeholder="https://example.com/image.jpg"
          required
        />
      </div>

      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <Label>Transformations</Label>
          <Select onValueChange={(v) => addTransformation(v as any)}>
            <SelectTrigger className="w-[200px]">
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

        {transformations.map((transform, index) => (
          <Card key={index} className="p-4">
            <div className="flex items-start justify-between gap-4">
              <div className="flex-1 space-y-3">
                <div className="font-medium capitalize">{transform.name}</div>
                {transform.name === 'resize' && (
                  <div className="grid grid-cols-2 gap-3">
                    <div>
                      <Label>Width</Label>
                      <Input
                        type="number"
                        value={transform.config.width}
                        onChange={(e) => updateTransformation(index, { ...transform.config, width: +e.target.value })}
                        min={1}
                      />
                    </div>
                    <div>
                      <Label>Height</Label>
                      <Input
                        type="number"
                        value={transform.config.height}
                        onChange={(e) => updateTransformation(index, { ...transform.config, height: +e.target.value })}
                        min={1}
                      />
                    </div>
                  </div>
                )}
                {transform.name === 'trim' && (
                  <div>
                    <Label>Threshold (0-255)</Label>
                    <Input
                      type="number"
                      value={transform.config.threshold}
                      onChange={(e) => updateTransformation(index, { threshold: +e.target.value })}
                      min={0}
                      max={255}
                    />
                  </div>
                )}
                {transform.name === 'blur' && (
                  <div>
                    <Label>Sigma</Label>
                    <Input
                      type="number"
                      step="0.1"
                      value={transform.config.sigma}
                      onChange={(e) => updateTransformation(index, { sigma: +e.target.value })}
                      min={0.1}
                    />
                  </div>
                )}
                {transform.name === 'rotate' && (
                  <div>
                    <Label>Angle</Label>
                    <Select
                      value={transform.config.angle.toString()}
                      onValueChange={(v) => updateTransformation(index, { angle: +v })}
                    >
                      <SelectTrigger>
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
              </div>
              <Button type="button" variant="ghost" size="icon" onClick={() => removeTransformation(index)}>
                <Trash2 className="w-4 h-4" />
              </Button>
            </div>
          </Card>
        ))}
      </div>

      <div className="flex gap-3">
        <Button type="submit" disabled={loading || !imageUrl || transformations.length === 0}>
          {loading ? 'Creating...' : 'Create Image'}
        </Button>
        {onCancel && (
          <Button type="button" variant="outline" onClick={onCancel}>
            Cancel
          </Button>
        )}
      </div>
    </form>
  );
}
