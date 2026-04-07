import React, { useState, useEffect } from 'react'
import api from '../api/client'

function Dashboard() {
  const [stats, setStats] = useState({ total_documents: 0, success: 0, failed: 0 })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadStats()
  }, [])

  const loadStats = async () => {
    try {
      const res = await api.get('/stats')
      setStats(res.data)
    } catch (err) {
      console.error('Failed to load stats:', err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-gray-800">Dashboard</h1>

      {loading ? (
        <div className="text-center py-12 text-gray-500">Loading...</div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <StatCard title="Total Documents" value={stats.total_documents} color="blue" />
          <StatCard title="Success" value={stats.success} color="green" />
          <StatCard title="Failed" value={stats.failed} color="red" />
        </div>
      )}

      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-xl font-semibold mb-4">Quick Actions</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <ActionCard title="Scan Document" href="/ocr" emoji="📷" />
          <ActionCard title="Manage Rooms" href="/rooms" emoji="🏠" />
          <ActionCard title="View Bookings" href="/bookings" emoji="📅" />
          <ActionCard title="Settings" href="/settings" emoji="⚙️" />
        </div>
      </div>
    </div>
  )
}

function StatCard({ title, value, color }) {
  const colors = {
    blue: 'bg-blue-100 text-blue-800',
    green: 'bg-green-100 text-green-800',
    red: 'bg-red-100 text-red-800',
  }

  return (
    <div className={`rounded-lg shadow p-6 ${colors[color]}`}>
      <div className="text-sm font-medium opacity-75">{title}</div>
      <div className="text-3xl font-bold mt-2">{value}</div>
    </div>
  )
}

function ActionCard({ title, href, emoji }) {
  return (
    <a href={href} className="block bg-gray-50 hover:bg-gray-100 rounded-lg p-4 text-center transition">
      <div className="text-3xl mb-2">{emoji}</div>
      <div className="font-medium text-gray-700">{title}</div>
    </a>
  )
}

export default Dashboard
