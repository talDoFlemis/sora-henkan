import { Button } from "@/components/ui/button"
import { useNavigate } from "react-router-dom"
import { Sparkles, Image as ImageIcon, Zap, ArrowRight } from "lucide-react"

export function LandingPage() {
  const navigate = useNavigate()

  return (
    <div className="min-h-screen w-full  flex items-center justify-center">
      <div className="">
        <div className="">
          {/* Hero Section */}
          <div className="text-center space-y-8 mb-20">
            <div className="inline-block">
              <div className="flex items-center gap-2 bg-white/80 backdrop-blur-sm px-4 py-2 rounded-full shadow-sm border border-gray-200">
                <Sparkles className="w-4 h-4 text-indigo-600" />
                <span className="text-sm font-medium text-gray-700">
                  AI-Powered Image Processing
                </span>
              </div>
            </div>

            <h1 className="text-7xl md:text-8xl font-bold tracking-tight">
              <span className="bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-600 bg-clip-text text-transparent">
                Sora Henkan
              </span>
            </h1>

            <p className="text-xl md:text-2xl text-gray-600 max-w-3xl mx-auto font-light">
              Transform your images instantly with powerful AI-driven
              processing.
              <br />
              Simple, fast, and beautiful.
            </p>

            <div className="flex flex-col sm:flex-row items-center justify-center gap-4 pt-4">
              <Button
                size="lg"
                onClick={() => navigate("/gallery")}
                className="text-lg px-10 py-7 rounded-full bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 shadow-lg hover:shadow-xl transition-all duration-200 group"
              >
                Get Started
                <ArrowRight className="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform" />
              </Button>
              <Button
                size="lg"
                onClick={() => navigate("/gallery")}
                className="text-lg px-10 py-7 rounded-full bg-gradient-to-r text-white from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 shadow-lg hover:shadow-xl transition-all duration-200 group"
              >
                View Gallery
              </Button>
            </div>
          </div>

          {/* Features Grid */}
          <div className="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto">
            <div className="group bg-white/60 backdrop-blur-sm p-8 rounded-3xl shadow-sm hover:shadow-lg transition-all duration-300 border border-white/60 hover:border-indigo-200">
              <div className="w-14 h-14 bg-gradient-to-br from-purple-500 to-pink-500 rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
                <Sparkles className="w-7 h-7 text-white" />
              </div>
              <h3 className="font-semibold text-2xl mb-3 text-gray-900">
                Multiple Transformations
              </h3>
              <p className="text-gray-600 leading-relaxed">
                Apply resize, blur, rotation, grayscale, and more with intuitive
                controls
              </p>
            </div>

            <div className="group bg-white/60 backdrop-blur-sm p-8 rounded-3xl shadow-sm hover:shadow-lg transition-all duration-300 border border-white/60 hover:border-indigo-200">
              <div className="w-14 h-14 bg-gradient-to-br from-blue-500 to-cyan-500 rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
                <Zap className="w-7 h-7 text-white" />
              </div>
              <h3 className="font-semibold text-2xl mb-3 text-gray-900">
                Real-time Updates
              </h3>
              <p className="text-gray-600 leading-relaxed">
                Watch your images transform instantly with live server-sent
                events
              </p>
            </div>

            <div className="group bg-white/60 backdrop-blur-sm p-8 rounded-3xl shadow-sm hover:shadow-lg transition-all duration-300 border border-white/60 hover:border-indigo-200">
              <div className="w-14 h-14 bg-gradient-to-br from-indigo-500 to-purple-500 rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
                <ImageIcon className="w-7 h-7 text-white" />
              </div>
              <h3 className="font-semibold text-2xl mb-3 text-gray-900">
                Gallery Management
              </h3>
              <p className="text-gray-600 leading-relaxed">
                Organize and manage all your transformed images in one beautiful
                place
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
