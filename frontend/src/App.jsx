import React from 'react'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import Dashboard from './pages/Dashboard'
import OCR from './pages/OCR'
import Rooms from './pages/Rooms'
import Bookings from './pages/Bookings'
import Settings from './pages/Settings'

function App() {
  return (
    <BrowserRouter>
      <div className="min-h-screen bg-gray-50">
        <nav className="bg-blue-600 text-white p-4">
          <div className="container mx-auto flex gap-6">
            <a href="/" className="font-bold text-lg">🏨 Hotel OCR</a>
            <a href="/" className="hover:text-blue-200">Dashboard</a>
            <a href="/ocr" className="hover:text-blue-200">OCR Scan</a>
            <a href="/rooms" className="hover:text-blue-200">Rooms</a>
            <a href="/bookings" className="hover:text-blue-200">Bookings</a>
            <a href="/settings" className="hover:text-blue-200">Settings</a>
          </div>
        </nav>
        <main className="container mx-auto p-6">
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/ocr" element={<OCR />} />
            <Route path="/rooms" element={<Rooms />} />
            <Route path="/bookings" element={<Bookings />} />
            <Route path="/settings" element={<Settings />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  )
}

export default App
