import { useEffect, useState } from 'react';
import { api, type Image } from '@/lib/api';
import { X, ImageIcon } from 'lucide-react';

export function AnnouncementBar() {
  const [latestImage, setLatestImage] = useState<Image | null>(null);
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    const eventSource = api.streamAllImages((image) => {
      setLatestImage(image);
      setVisible(true);

      // Auto hide after 10 seconds
      setTimeout(() => setVisible(false), 10000);
    });

    return () => eventSource.close();
  }, []);

  if (!visible || !latestImage) return null;

  return (
    <div className="fixed top-0 left-0 right-0 z-50 bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-600 text-white py-3 px-4 shadow-lg animate-slide-down">
      <div className="container mx-auto flex items-center justify-between">
        <div className="flex items-center gap-3 flex-1 overflow-hidden">
          <ImageIcon className="w-5 h-5 flex-shrink-0 animate-pulse" />
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium truncate">
              New image processed:{' '}
              <span className="font-mono">
                {latestImage.id.slice(0, 8)}...
              </span>
            </p>
            <p className="text-xs opacity-90">
              Status: {latestImage.status} •{' '}
              {latestImage.transformations.length} transformation(s)
            </p>
          </div>
        </div>
        <button
          onClick={() => setVisible(false)}
          className="ml-4 hover:bg-white/20 rounded-full p-1.5 transition-colors flex-shrink-0"
          aria-label="Close"
        >
          <X className="w-4 h-4" />
        </button>
      </div>
    </div>
  );
}
