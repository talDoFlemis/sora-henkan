import { Button } from "@/components/ui/button"
import { useNavigate } from "react-router-dom"
import {
  Sparkles,
  Image as ImageIcon,
  Zap,
  ArrowRight,
  Layers,
  Wand2,
  Camera,
} from "lucide-react"

export function LandingPage() {
  const navigate = useNavigate()

  return (
    <div className="min-h-screen w-full relative overflow-hidden">
      {/* Animated background elements */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-purple-300 rounded-full mix-blend-multiply filter blur-3xl opacity-30 animate-float" />
        <div
          className="absolute top-40 -left-40 w-80 h-80 bg-indigo-300 rounded-full mix-blend-multiply filter blur-3xl opacity-30 animate-float"
          style={{ animationDelay: "1s" }}
        />
        <div
          className="absolute bottom-40 right-40 w-80 h-80 bg-pink-300 rounded-full mix-blend-multiply filter blur-3xl opacity-30 animate-float"
          style={{ animationDelay: "2s" }}
        />
      </div>

      <div className="relative z-10 w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Hero Section */}
        <div className="flex flex-col items-center text-center space-y-8 pt-16 pb-24 animate-fade-in-up">
          {/* Badge */}
          <div className="inline-flex items-center gap-2 px-5 py-2.5 rounded-full glass shadow-lg border border-white/50">
            <Sparkles className="w-4 h-4 text-indigo-600" />
            <span className="text-sm font-semibold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
              Powered by Advanced AI
            </span>
          </div>

          {/* Main Title */}
          <div className="space-y-4">
            <h1 className="text-6xl md:text-7xl lg:text-8xl font-black tracking-tight">
              <span className="gradient-text">Sora</span>
              <span className="text-gray-800"> Henkan</span>
            </h1>
            <p className="text-xl md:text-2xl text-gray-600 max-w-2xl mx-auto font-light leading-relaxed">
              Transform your images with powerful AI-driven processing.
              <br />
              <span className="font-medium text-gray-700">
                Simple. Fast. Beautiful.
              </span>
            </p>
          </div>

          {/* CTA Buttons */}
          <div className="flex flex-col sm:flex-row items-center gap-4 pt-6">
            <Button
              size="lg"
              onClick={() => navigate("/gallery")}
              className="group relative px-8 py-6 text-lg font-semibold rounded-2xl gradient-primary text-white shadow-xl hover:shadow-2xl transition-all duration-300 hover:-translate-y-1 overflow-hidden"
            >
              <span className="relative z-10 flex items-center gap-2">
                Get Started Free
                <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
              </span>
              <div className="absolute inset-0 bg-white/20 translate-y-full group-hover:translate-y-0 transition-transform duration-300" />
            </Button>
            <Button
              size="lg"
              variant="outline"
              onClick={() => navigate("/gallery")}
              className="px-8 py-6 text-lg font-semibold rounded-2xl glass border-2 border-indigo-200 hover:border-indigo-400 hover:bg-white/80 transition-all duration-300 text-gray-700"
            >
              <Camera className="w-5 h-5 mr-2" />
              View Gallery
            </Button>
          </div>

          {/* Stats */}
          <div className="flex flex-wrap justify-center gap-8 pt-8">
            <div className="text-center">
              <div className="text-3xl font-bold gradient-text">10+</div>
              <div className="text-sm text-gray-500">Transformations</div>
            </div>
            <div className="text-center">
              <div className="text-3xl font-bold gradient-text">&lt;2s</div>
              <div className="text-sm text-gray-500">Processing Time</div>
            </div>
            <div className="text-center">
              <div className="text-3xl font-bold gradient-text">∞</div>
              <div className="text-sm text-gray-500">Possibilities</div>
            </div>
          </div>
        </div>

        {/* Features Section */}
        <div className="py-16">
          <div className="text-center mb-12">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-800 mb-4">
              Everything You Need
            </h2>
            <p className="text-gray-600 max-w-xl mx-auto">
              Professional image transformation tools at your fingertips
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto stagger-children">
            {/* Feature 1 */}
            <div className="group p-8 rounded-3xl glass-dark shadow-lg hover-lift cursor-pointer">
              <div className="w-16 h-16 rounded-2xl gradient-primary flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300 shadow-lg">
                <Wand2 className="w-8 h-8 text-white" />
              </div>
              <h3 className="text-xl font-bold text-gray-800 mb-3">
                Multiple Transformations
              </h3>
              <p className="text-gray-600 leading-relaxed">
                Resize, blur, rotate, grayscale, and trim with intuitive
                controls
              </p>
            </div>

            {/* Feature 2 */}
            <div className="group p-8 rounded-3xl glass-dark shadow-lg hover-lift cursor-pointer">
              <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-cyan-500 to-blue-600 flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300 shadow-lg">
                <Zap className="w-8 h-8 text-white" />
              </div>
              <h3 className="text-xl font-bold text-gray-800 mb-3">
                Real-time Updates
              </h3>
              <p className="text-gray-600 leading-relaxed">
                Watch your images transform instantly with live server-sent
                events
              </p>
            </div>

            {/* Feature 3 */}
            <div className="group p-8 rounded-3xl glass-dark shadow-lg hover-lift cursor-pointer">
              <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-pink-500 to-rose-600 flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300 shadow-lg">
                <Layers className="w-8 h-8 text-white" />
              </div>
              <h3 className="text-xl font-bold text-gray-800 mb-3">
                Gallery Management
              </h3>
              <p className="text-gray-600 leading-relaxed">
                Organize and manage all your transformed images in one beautiful
                place
              </p>
            </div>
          </div>
        </div>

        {/* How it works */}
        <div className="py-16">
          <div className="text-center mb-12">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-800 mb-4">
              How It Works
            </h2>
            <p className="text-gray-600 max-w-xl mx-auto">
              Three simple steps to transform your images
            </p>
          </div>

          <div className="flex flex-col md:flex-row items-center justify-center gap-8 max-w-4xl mx-auto">
            {/* Step 1 */}
            <div className="flex flex-col items-center text-center flex-1">
              <div className="w-20 h-20 rounded-full gradient-primary flex items-center justify-center text-white text-2xl font-bold mb-4 shadow-xl">
                1
              </div>
              <h3 className="font-bold text-gray-800 mb-2">Upload Image</h3>
              <p className="text-sm text-gray-600">Paste any image URL</p>
            </div>

            <div className="hidden md:block w-24 h-0.5 bg-gradient-to-r from-indigo-300 to-purple-300" />

            {/* Step 2 */}
            <div className="flex flex-col items-center text-center flex-1">
              <div className="w-20 h-20 rounded-full bg-gradient-to-br from-cyan-500 to-blue-600 flex items-center justify-center text-white text-2xl font-bold mb-4 shadow-xl">
                2
              </div>
              <h3 className="font-bold text-gray-800 mb-2">
                Choose Transforms
              </h3>
              <p className="text-sm text-gray-600">Select effects to apply</p>
            </div>

            <div className="hidden md:block w-24 h-0.5 bg-gradient-to-r from-blue-300 to-pink-300" />

            {/* Step 3 */}
            <div className="flex flex-col items-center text-center flex-1">
              <div className="w-20 h-20 rounded-full bg-gradient-to-br from-pink-500 to-rose-600 flex items-center justify-center text-white text-2xl font-bold mb-4 shadow-xl">
                3
              </div>
              <h3 className="font-bold text-gray-800 mb-2">Get Results</h3>
              <p className="text-sm text-gray-600">Download or share</p>
            </div>
          </div>
        </div>

        {/* CTA Section */}
        <div className="py-16">
          <div className="relative rounded-3xl overflow-hidden p-12 text-center gradient-primary shadow-2xl">
            <div className="absolute inset-0 bg-black/10" />
            <div className="relative z-10">
              <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">
                Ready to Transform?
              </h2>
              <p className="text-white/80 mb-8 max-w-xl mx-auto">
                Start transforming your images today. No account required.
              </p>
              <Button
                size="lg"
                onClick={() => navigate("/gallery")}
                className="px-10 py-6 text-lg font-semibold rounded-2xl bg-white text-indigo-600 hover:bg-gray-100 shadow-xl hover:shadow-2xl transition-all duration-300"
              >
                <ImageIcon className="w-5 h-5 mr-2" />
                Open Gallery
              </Button>
            </div>
          </div>
        </div>

        {/* Footer */}
        <footer className="py-8 text-center text-gray-500 text-sm">
          <p>© 2025 Sora Henkan. Built with ❤️ using Go, React & DynamoDB</p>
        </footer>
      </div>
    </div>
  )
}
