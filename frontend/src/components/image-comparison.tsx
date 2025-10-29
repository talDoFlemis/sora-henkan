import { useState } from "react"
import { Card } from "@/components/ui/card"

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

  const handleMove = (
    e: React.MouseEvent<HTMLDivElement> | React.TouchEvent<HTMLDivElement>,
  ) => {
    const rect = e.currentTarget.getBoundingClientRect()
    const x = "touches" in e ? e.touches[0].clientX : e.clientX
    const position = ((x - rect.left) / rect.width) * 100
    setSliderPosition(Math.max(0, Math.min(100, position)))
  }

  return (
    <Card
      className="relative overflow-hidden select-none"
      onMouseMove={handleMove}
      onTouchMove={handleMove}
    >
      <div className="relative w-full aspect-video">
        <img
          src={beforeSrc}
          alt={`${alt} - before`}
          className="absolute inset-0 w-full h-full object-contain"
        />
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
        <div
          className="absolute inset-y-0 w-1 bg-white cursor-ew-resize"
          style={{ left: `${sliderPosition}%` }}
        >
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-8 h-8 bg-white rounded-full shadow-lg flex items-center justify-center">
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M8 9l4-4 4 4m0 6l-4 4-4-4"
              />
            </svg>
          </div>
        </div>
      </div>
    </Card>
  )
}
