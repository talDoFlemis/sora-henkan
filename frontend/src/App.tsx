import { BrowserRouter, Routes, Route } from "react-router-dom"
import "./App.css"
import { LandingPage } from "./pages/landing"
import { GalleryPage } from "./pages/gallery"
import { ImageDetailPage } from "./pages/image-detail"

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/gallery" element={<GalleryPage />} />
        <Route path="/images/:id" element={<ImageDetailPage />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
