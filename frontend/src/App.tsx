import { useState } from 'react'
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import { env } from './utils/constants'
import { LandingPage } from './pages/landing';
import { GalleryPage } from './pages/gallery';
import { ImageDetailPage } from './pages/image-detail';
import { AnnouncementBar } from './components/announcement-bar';

function App() {
  const [count, setCount] = useState(0)

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
