import { useState, useEffect } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

function App() {
  const [file, setFile] = useState(null);
  const [audioURL, setAudioURL] = useState(null);
  const [loading, setLoading] = useState(null);
  const [bucketFiles, setBucketFiles] = useState([]);
  const [selectedFile, setSelectedFile] = useState(null);

  const updateSelectedFile = (newFile) => {
    setSelectedFile(newFile);
  }
  
  const handleFileChange = (event) => {
    setFile(event.target.files[0]);
  }

  useEffect(() => {
    async function fetchFiles() {
      try {
        const res = await fetch("http://localhost:3000/bucket-files");
        const data = await res.json();
        console.log("Data from frontend: ", data);
        setBucketFiles(data);
      } catch (err) {
        console.log("Error grabbing buckets from s3 bucket", err);
      }
    }
    fetchFiles();
  }, []);

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
      <h3>Bucket Files</h3>
      <ul style={{ listStyle: 'none', padding: 0}}>
      {bucketFiles.map((file) => (
        <li 
          key={file.id}
          onClick={() => updateSelectedFile(file)}
          style={{
            cursor: 'pointer',
            padding: '0.5rem',
            border: '1px solid #ccc',
            borderRadius: '8px',
            marginBottom: '0.5rem',
            backgroundColor:
            selectedFile?.id === file.id ? '#d1e7ff' : 'white',
            transition: 'background-color 0.15s ease',
          }}
        >
          {file.title}
        </li>
      ))}
    </ul>
    </div>
  )
}

export default App
