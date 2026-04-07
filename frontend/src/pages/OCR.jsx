import React, { useState, useRef } from 'react'
import api from '../api/client'

function OCR() {
  const [file, setFile] = useState(null)
  const [preview, setPreview] = useState(null)
  const [result, setResult] = useState(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const fileInputRef = useRef(null)

  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0]
    if (selectedFile) {
      setFile(selectedFile)
      setPreview(URL.createObjectURL(selectedFile))
      setResult(null)
      setError(null)
    }
  }

  const handleDrop = (e) => {
    e.preventDefault()
    const droppedFile = e.dataTransfer.files[0]
    if (droppedFile && droppedFile.type.startsWith('image/')) {
      setFile(droppedFile)
      setPreview(URL.createObjectURL(droppedFile))
      setResult(null)
      setError(null)
    }
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (!file) return

    setLoading(true)
    setError(null)

    const formData = new FormData()
    formData.append('image', file)

    try {
      const res = await api.post('/documents/process', formData)
      setResult(res.data)
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to process image')
    } finally {
      setLoading(false)
    }
  }

  const resetForm = () => {
    setFile(null)
    setPreview(null)
    setResult(null)
    setError(null)
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-gray-800">📷 OCR Scanner</h1>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold mb-4">Upload Document</h2>

          <div
            className="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center cursor-pointer hover:border-blue-500 transition"
            onDrop={handleDrop}
            onDragOver={(e) => e.preventDefault()}
            onClick={() => fileInputRef.current?.click()}
          >
            <input
              ref={fileInputRef}
              type="file"
              accept="image/*"
              onChange={handleFileChange}
              className="hidden"
            />
            <div className="text-6xl mb-4">📄</div>
            <p className="text-gray-600">
              Drag & drop image here or <span className="text-blue-600">browse</span>
            </p>
            <p className="text-sm text-gray-400 mt-2">Supports: JPG, PNG, PDF</p>
          </div>

          {preview && (
            <div className="mt-4">
              <img src={preview} alt="Preview" className="max-h-64 mx-auto rounded-lg shadow" />
            </div>
          )}

          {file && (
            <div className="mt-4 flex gap-4">
              <button
                onClick={handleSubmit}
                disabled={loading}
                className="flex-1 bg-blue-600 text-white py-3 px-6 rounded-lg font-medium hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition"
              >
                {loading ? 'Processing...' : 'Process Image'}
              </button>
              <button
                onClick={resetForm}
                className="px-6 py-3 border border-gray-300 rounded-lg font-medium hover:bg-gray-50 transition"
              >
                Clear
              </button>
            </div>
          )}

          {error && (
            <div className="mt-4 p-4 bg-red-100 text-red-700 rounded-lg">
              ❌ {error}
            </div>
          )}
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold mb-4">Result</h2>

          {result ? (
            <div className="space-y-4">
              <ResultField label="Guest Name" value={result.guest_name} />
              <ResultField label="ID Card" value={result.id_card} />
              <ResultField label="Phone" value={result.phone} />
              <ResultField label="Room Number" value={result.room_number} />
              <ResultField label="Check-in Date" value={result.check_in_date} />
              <ResultField label="Check-out Date" value={result.check_out_date} />
              <ResultField label="License Plate" value={result.license_plate} />
              <div className="flex items-center gap-2">
                <span className="font-medium">Confidence:</span>
                <div className="flex-1 bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-green-500 h-2 rounded-full"
                    style={{ width: `${(result.confidence || 0) * 100}%` }}
                  />
                </div>
                <span>{((result.confidence || 0) * 100).toFixed(0)}%</span>
              </div>
              <div className="pt-4 border-t">
                <p className="text-sm text-gray-500">Provider: {result.provider}</p>
                <p className="text-sm text-gray-500">Document ID: {result.id}</p>
              </div>
            </div>
          ) : (
            <div className="text-center py-12 text-gray-400">
              <div className="text-6xl mb-4">📋</div>
              <p>Upload an image to see the result</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

function ResultField({ label, value }) {
  return (
    <div className="flex gap-4">
      <span className="font-medium text-gray-600 w-32">{label}:</span>
      <span className="text-gray-800">{value || '-'}</span>
    </div>
  )
}

export default OCR
