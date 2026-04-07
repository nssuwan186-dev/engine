import React, { useState } from 'react'

function Settings() {
  const [settings, setSettings] = useState({
    apiUrl: localStorage.getItem('api_url') || 'http://localhost:8080/api/v1',
    theme: localStorage.getItem('theme') || 'light',
  })
  const [saved, setSaved] = useState(false)

  const handleSave = () => {
    localStorage.setItem('api_url', settings.apiUrl)
    localStorage.setItem('theme', settings.theme)
    setSaved(true)
    setTimeout(() => setSaved(false), 2000)
  }

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-gray-800">⚙️ Settings</h1>

      <div className="bg-white rounded-lg shadow p-6 space-y-6">
        <div>
          <h2 className="text-lg font-semibold mb-4">API Configuration</h2>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            API URL
          </label>
          <input
            type="text"
            value={settings.apiUrl}
            onChange={(e) => setSettings({ ...settings, apiUrl: e.target.value })}
            className="w-full border rounded-lg px-4 py-2"
            placeholder="http://localhost:8080/api/v1"
          />
        </div>

        <div>
          <h2 className="text-lg font-semibold mb-4">Appearance</h2>
          <select
            value={settings.theme}
            onChange={(e) => setSettings({ ...settings, theme: e.target.value })}
            className="w-full border rounded-lg px-4 py-2"
          >
            <option value="light">Light</option>
            <option value="dark">Dark</option>
          </select>
        </div>

        <div>
          <h2 className="text-lg font-semibold mb-4">Admin Access</h2>
          <p className="text-sm text-gray-500 mb-4">
            Use CLI backdoor for administrative tasks.
          </p>
          <code className="block bg-gray-100 p-4 rounded-lg text-sm">
            ./hotel-ocr-cli -cmd=shell -keyfile=./data/backdoor.key
          </code>
        </div>

        <div className="pt-4">
          <button
            onClick={handleSave}
            className="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700"
          >
            Save Settings
          </button>
          {saved && <span className="ml-4 text-green-600">✓ Saved!</span>}
        </div>
      </div>
    </div>
  )
}

export default Settings
