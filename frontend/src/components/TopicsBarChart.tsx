import { useEffect, useState } from 'react'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'

interface Analytics {
  total_submissions: number
  avg_completion: number
  topics: Record<string, number>
}

export default function TopicsBarChart({ analytics }: { analytics?: Analytics }) {
  const [chartData, setChartData] = useState<{ topic: string; count: number }[]>([])

  useEffect(() => {
    if (analytics?.topics) {
      // Convert topics object into array format for Recharts
      const formattedData = Object.entries(analytics.topics).map(([topic, count]) => ({
        topic,
        count
      }))
      setChartData(formattedData)
    }
  }, [analytics])

  return (
    <div style={{ width: '100%', height: '100%' }}>
      <ResponsiveContainer width="100%" height="100%">
        <BarChart
          data={chartData}
          margin={{ top: 20, right: 20, left: 0, bottom: 20 }}
        >
          <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
          <XAxis dataKey="topic" angle={-30} textAnchor="end" height={70} tick={{ fill: '#374151' }} />
          <YAxis tick={{ fill: '#374151' }} />
          <Tooltip cursor={{ fill: 'rgba(0,0,0,0.05)' }} />
          <Bar dataKey="count" fill="#6366f1" radius={[6, 6, 0, 0]} />
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}
