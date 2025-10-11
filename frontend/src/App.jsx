import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

function App() {
  const [file, setFile] = useState(null);
  const [audioURL, setAudioURL] = useState(null);
  const [loading, setLoading] = useState(null);

  
  const handleFileChange = (event) => {
    setFile(event.target.files[0]);
  }

  return (
    <div>
      <input type="file" onChange={handleFileChange} />
      {file && <p>Selected file: {file.name}</p>}
      {file && (
        <button onClick={()=> setAudioURL("/example.mp3")}>
          Narrate
        </button>
      )}
      {audioURL && (
        <div>
          <h3>Your Narration:</h3>
            <audio controls src={audioURL} />
        </div>
      )}
    </div>
  )
}

export default App
