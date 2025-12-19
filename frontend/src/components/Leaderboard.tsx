import { useEffect, useState } from 'react'
import '../App.css'
import { motion, AnimatePresence } from 'framer-motion'

interface Entry {
  Member: string;
  Score: number;
}

const Leaderboard = () => {
  const [connectionString, setConnectionString] = useState('Not Connected')
  const [leaderboardEntries, setLeaderboardEntries] = useState<Entry[]>([])

  useEffect(() => {
    // Note: When you deploy, change 'localhost' to your Fly.io/Render URL
    const ws = new WebSocket("ws://localhost:8080/ws")

    ws.onopen = () => {
      setConnectionString("Connected")
    }

    ws.onclose = () => {
      setConnectionString("Disconnected")
    }

    ws.onmessage = (event) => {
      try {
        // Your Go server sends: [{username: "user1", score: 100}, ...]
        const data = JSON.parse(event.data)
        
        // Directly set the data because Go sends the array as the top-level object
        setLeaderboardEntries(data || [])
      } catch (err) {
        console.error("Error parsing leaderboard JSON:", err)
      }
    }

    return () => ws.close()
  }, [])

  return (
    <div className="leaderboard-page">
      <div className="status-bar">
         <span className={`dot ${connectionString === 'Connected' ? 'online' : 'offline'}`}></span>
         {connectionString}
      </div>

      <div className="leaderboard-card">
        <AnimatePresence mode="popLayout">
          {leaderboardEntries.length > 0 ? (
            <motion.ul
              layout
              className="leaderboard-list"
            >
              {leaderboardEntries.map((entry, index) => (
                <motion.li
                  layout
                  // IMPORTANT: Using username as key allows Framer Motion 
                  // to animate the row moving UP or DOWN when ranks change
                  key={entry.Member} 
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, scale: 0.95 }}
                  transition={{ 
                    layout: { type: "spring", stiffness: 300, damping: 30 },
                    opacity: { duration: 0.2 } 
                  }}
                  className="leaderboard-row"
                >
                  <div className="player-info">
                    <span className="rank">{index + 1}</span>
                    <span className="username">{entry.Member}</span>
                  </div>
                  
                  <motion.span
                    // This creates a "pulse" effect whenever the score updates
                    key={entry.Score}
                    initial={{ scale: 1.2, color: "#f39c12" }}
                    animate={{ scale: 1, color: "#000" }}
                    className="score"
                  >
                    {entry.Score.toLocaleString()}
                  </motion.span>
                </motion.li>
              ))}
            </motion.ul>
          ) : (
            <div className="no-data">Waiting for scores...</div>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
}

export default Leaderboard;