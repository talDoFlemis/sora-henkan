import { useState } from "react"

interface ImageComparisonProps {
  beforeSrc: string
  afterSrc: string
  alt?: string
}

export function ImageComparison({
  beforeSrc,
  afterSrc,
  alt = "Image comparison",
}: ImageComparisonProps) {
  const [sliderPosition, setSliderPosition] = useState(50)
  const [isDragging, setIsDragging] = useState(false)

  const handleMove = (
    e: React.MouseEvent<HTMLDivElement> | React.TouchEvent<HTMLDivElement>,
  ) => {
    if (!isDragging && e.type !== "touchmove") return
    const rect = e.currentTarget.getBoundingClientRect()
    const x = "touches" in e ? e.touches[0].clientX : e.clientX
    const position = ((x - rect.left) / rect.width) * 100
    setSliderPosition(Math.max(0, Math.min(100, position)))
  }

  return (
    <div
      className="relative overflow-hidden select-none cursor-ew-resize"
      onMouseDown={() => setIsDragging(true)}
      onMouseUp={() => setIsDragging(false)}
      onMouseLeave={() => setIsDragging(false)}
      onMouseMove={handleMove}
      onTouchMove={handleMove}
      onTouchStart={() => setIsDragging(true)}
      onTouchEnd={() => setIsDragging(false)}
    >
      {/* Labels */}
      <div className="absolute top-4 left-4 z-20">
        <span className="px-3 py-1.5 bg-black/60 backdrop-blur-sm text-white text-xs font-medium rounded-full">
          Original
        </span>
      </div>
      <div className="absolute top-4 right-4 z-20">
        <span className="px-3 py-1.5 bg-gradient-to-r from-indigo-600 to-purple-600 text-white text-xs font-medium rounded-full">
          Transformed
        </span>
      </div>

      <div className="relative w-full aspect-video bg-gray-100">
        {/* Before Image */}
        <img
          src={beforeSrc}
          alt={`${alt} - before`}
          className="absolute inset-0 w-full h-full object-contain"
        />
        
        {/* After Image with clip */}
        <div
          className="absolute inset-0 overflow-hidden"
          style={{ clipPath: `inset(0 ${100 - sliderPosition}% 0 0)` }}
        >
          <img
            src={afterSrc}
            alt={`${alt} - after`}
            className="absolute inset-0 w-full h-full object-contain"
          />
        </div>

        {/* Slider Line */}
        <div
          className="absolute inset-y-0 w-0.5 bg-white shadow-lg z-10"
          style={{ left: `${sliderPosition}%`, transform: "translateX(-50%)" }}
        >
          {/* Slider Handle */}
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-10 h-10 bg-white rounded-full shadow-xl flex items-center justify-center border-2 border-indigo-500 transition-transform hover:scale-110">
            <svg
              className="w-5 h-5 text-indigo-600"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2.5}
                d="M8 7l-4 5 4 5M16 7l4 5-4 5"
              />
            </svg>
          </div>
        </div>
      </div>

      {/* Hint */}
      <div className="absolute bottom-4 left-1/2 -translate-x-1/2 z-20">
        <span className="px-3 py-1.5 bg-black/40 backdrop-blur-sm text-white/80 text-xs rounded-full">
          Drag to compare
        </span>
      </div>
    </div>
  )
}
