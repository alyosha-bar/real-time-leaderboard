import { useEffect, useState } from 'react'
import '../App.css'
import { motion, AnimatePresence } from 'framer-motion'





interface Entry {
  username: string,
  score: number
}

const Leaderboard = () => {

    const [connectionString, setConnectionString] = useState('Not Connected to Websockets.')
    const [leaderboardEntries, setLeaderboardEntries] = useState<Entry[]>()


    useEffect(() => {
        const ws = new WebSocket("ws://localhost:8080/ws")

        ws.onopen = () => {
            setConnectionString("✅ Connected to leaderboard updates.")
        }

        ws.onclose = () => {
            setConnectionString("❌ Disconnected from WebSockets.")
        }

        ws.onerror = (err) => {
            console.error("WebSocket error:", err)
            setConnectionString("⚠️ WebSocket error - check console")
        }

        ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data)
                setLeaderboardEntries(data.entities || [])
            } catch (err) {
                console.error("Error parsing message:", err)
            }
        }

        // cleanup function
        return () => {
            ws.close()
        }
    }, [])


    return ( 
        <div className="leaderboard-page">
            {/* <h1 className="title"> Live Leaderboard </h1> */}

            <div className="leaderboard-card">
                {leaderboardEntries && (
                <AnimatePresence>
                    {leaderboardEntries.length > 0 ? (
                    <motion.ul
                        layout
                        className="leaderboard-list"
                        transition={{ layout: { duration: 0.4, type: 'spring' } }}
                    >
                        {leaderboardEntries.map((entry, index) => (
                        <motion.li
                            layout
                            key={entry.username}
                            initial={{ opacity: 0, y: -10 }}
                            animate={{ opacity: 1, y: 0 }}
                            exit={{ opacity: 0 }}
                            transition={{ duration: 0.3 }}
                            className="leaderboard-row"
                        >
                            <div className="player-info">
                            <motion.span layout className="rank">
                                {index + 1}.
                            </motion.span>
                            <span className="username">{entry.username}</span>
                            </div>
                            <motion.span
                            layout
                            animate={{ scale: [1, 1.1, 1] }}
                            transition={{ duration: 0.3 }}
                            className="score"
                            >
                            {entry.score}
                            </motion.span>
                        </motion.li>
                        ))}
                    </motion.ul>
                    ) : (
                    <div className="no-data">No leaderboard data</div>
                    )}
                </AnimatePresence>
                )}
            </div>

            {/* <div className="connection-status">{connectionString}</div> */}
            </div>
    );
}
 
export default Leaderboard;