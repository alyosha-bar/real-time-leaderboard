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


  return ( 
    <>
      <div className='dashboard'>
        {/* <div className='main'>
          <div className='metrics'>
            <span className='metric'>Total Submissions: {analytics?.total_submissions ?? 0}</span>
            <span className='metric'>Avg Completion: {analytics?.avg_completion ?? 0}s</span>
          </div>

          <div className='chart'>
            <TopicsBarChart analytics={analytics}></TopicsBarChart>
          </div>
        </div> */}

        <div className='leaderboard-container'>
          <Leaderboard />
        </div>
      </div>
    </>
  );
}

export default App;
