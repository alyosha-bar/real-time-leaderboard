import { useEffect, useState } from 'react';
import './App.css'
import Leaderboard from './components/Leaderboard';
import TopicsBarChart from './components/TopicsBarChart';

interface Analytics {
  total_submissions: number,
  avg_completion: number,
  topics: Record<string, number>
}


const App = () => {

  const [analytics, setAnalytics] = useState<Analytics>()

  useEffect( () => {
    const fetchAnalytics = async () => {

      try {
        const response = await fetch("http://localhost:8080/analytics")

        const data = await response.json()

        console.log(data)

        setAnalytics(data)
      } catch (err) {
        console.error(err)
      }
      
    }

    fetchAnalytics()
  }, [])

  return ( 
    <>
      <div className='dashboard'>
        <div className='main'>
          <div className='metrics'>
            <span className='metric'>Total Submissions: {analytics?.total_submissions ?? 0}</span>
            <span className='metric'>Avg Completion: {analytics?.avg_completion ?? 0}s</span>
            <span className='metric'>Another Metric</span>
          </div>

          <div className='chart'>
            <TopicsBarChart analytics={analytics}></TopicsBarChart>
          </div>
        </div>

        <div className='leaderboard-container'>
          <Leaderboard />
        </div>
      </div>
    </>
  );
}

export default App;
