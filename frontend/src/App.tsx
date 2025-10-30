import './App.css'
import Leaderboard from './components/Leaderboard';


const App = () => {
  return ( 
    <>
      <div className='dashboard'>
        <div className='main'>
          <div className='metrics'>
            <span className='metric'></span>  
            <span className='metric'></span>  
            <span className='metric'></span>  
          </div>
          <div className='chart'> </div>
        </div>
        <div className='leaderboard-container'>
          <Leaderboard />
        </div>
      </div>
    

    </>
  );
}

export default App;
