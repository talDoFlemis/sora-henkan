import { Button } from '@/components/ui/button';
import { useNavigate } from 'react-router-dom';
import { Sparkles, Image as ImageIcon, Zap } from 'lucide-react';

export function LandingPage() {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-purple-50 to-pink-50">
      <div className="container mx-auto px-4 py-16">
        <div className="max-w-4xl mx-auto text-center space-y-8">
          <div className="space-y-4">
            <h1 className="text-6xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
              Sora Henkan
            </h1>
            <p className="text-2xl text-gray-600">Transform Your Images with AI-Powered Magic</p>
          </div>

          <p className="text-lg text-gray-700 max-w-2xl mx-auto leading-relaxed">
            Sora Henkan is a powerful image processing service that allows you to transform, resize, enhance, 
            and manipulate images with ease. Upload your images and apply multiple transformations in real-time 
            with our intuitive interface.
          </p>

          <div className="grid md:grid-cols-3 gap-6 my-12">
            <div className="bg-white p-6 rounded-lg shadow-lg">
              <Sparkles className="w-12 h-12 mx-auto mb-4 text-purple-600" />
              <h3 className="font-semibold text-xl mb-2">Multiple Transformations</h3>
              <p className="text-gray-600">Apply resize, blur, rotation, grayscale, and more</p>
            </div>
            <div className="bg-white p-6 rounded-lg shadow-lg">
              <Zap className="w-12 h-12 mx-auto mb-4 text-blue-600" />
              <h3 className="font-semibold text-xl mb-2">Real-time Updates</h3>
              <p className="text-gray-600">Watch your images transform in real-time with SSE</p>
            </div>
            <div className="bg-white p-6 rounded-lg shadow-lg">
              <ImageIcon className="w-12 h-12 mx-auto mb-4 text-pink-600" />
              <h3 className="font-semibold text-xl mb-2">Gallery Management</h3>
              <p className="text-gray-600">Manage all your transformed images in one place</p>
            </div>
          </div>

          <Button size="lg" onClick={() => navigate('/gallery')} className="text-lg px-8 py-6">
            Get Started
          </Button>
        </div>
      </div>
    </div>
  );
}
