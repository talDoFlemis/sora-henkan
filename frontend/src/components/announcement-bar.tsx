import { useEffect, useState } from 'react';
import { api, type Image } from '@/lib/api';
import { X } from 'lucide-react';

export function AnnouncementBar() {
  const [updates, setUpdates] = useState<Image[]>([]);
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    const eventSource = api.streamAllImages((image) => {
      setUpdates((prev) => [image, ...prev.slice(0, 4)]);
      setVisible(true);
    });

    return () => eventSource.close();
  }, []);

  if (!visible || updates.length === 0) return null;

  return (
    <div className="fixed top-0 left-0 right-0 z-50 bg-gradient-to-r from-blue-600 to-purple-600 text-white py-2 px-4 shadow-lg">
      <div className="container mx-auto flex items-center justify-between">
        <div className="flex-1 overflow-hidden">
          <p className="text-sm font-medium animate-slide-left">
            🎉 New image processed: {updates[0]?.id} - Status: {updates[0]?.status}
          </p>
        </div>
        <button onClick={() => setVisible(false)} className="ml-4 hover:bg-white/20 rounded p-1">
          <X className="w-4 h-4" />
        </button>
      </div>
    </div>
  );
}
