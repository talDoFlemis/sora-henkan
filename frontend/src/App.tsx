import { BrowserRouter, Routes, Route } from "react-router-dom"
import "./App.css"
import { LandingPage } from "./pages/landing"
import { GalleryPage } from "./pages/gallery"
import { ImageDetailPage } from "./pages/image-detail"
import { AnnouncementBar } from "./components/announcement-bar"

function App() {
  return (
    <BrowserRouter>
      <AnnouncementBar />
      <div className="pt-10">
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route path="/gallery" element={<GalleryPage />} />
          <Route path="/images/:id" element={<ImageDetailPage />} />
        </Routes>
      </div>
    </BrowserRouter>
  )
}

export default App
