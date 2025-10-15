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
  const [uploading, setUploading] = useState(false);
  const [uploadSuccess, setUploadSuccess] = useState(false);

  const updateSelectedFile = (newFile) => {
    setSelectedFile(newFile);
  }
  
  const handleFileChange = (event) => {
    setFile(event.target.files[0]);
  }

  const handleUpload = async () => {
    if (!file) {
      return;
    }
    updateSelectedFile(file);
    setUploading(true);
    setUploadSuccess(false);

    try {
      const formData = new FormData();
      formData.append('file', file);
      const res = await fetch('http://localhost:3000/upload', {
        method: 'POST',
        body: formData,
      });
      if (!res.ok) {
        throw new Error('Upload failed!')
      }
      setUploadSuccess(true);
      setFile(null);
      const updated = await fetch('http://localhost:3000/bucket-files');
      const data = await updated.json();
      setBucketFiles(data);
    } catch (err) {
      console.error('Error uploading file', err);
      alert('Upload failed. Check console for details.');
    } finally {
      setUploading(false);
    }
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
    <div style={{ padding: '2rem' }}>
      <h2>Upload your note</h2>
      <input type="file" onChange={handleFileChange} />
      {file && <p>Selected: {file.name}</p>}

      {file && (
        <button
          onClick={handleUpload}
          disabled={uploading}
          style={{
            marginTop: '1rem',
            padding: '0.5rem 1rem',
            cursor: uploading ? 'not-allowed' : 'pointer',
          }}
        >
          {uploading ? 'Uploading...' : 'Upload'}
        </button>
      )}

      {uploadSuccess && (
        <p style={{ color: 'green', marginTop: '1rem' }}>
          ✅ File uploaded successfully!
        </p>
      )}

      <h3 style={{ marginTop: '2rem' }}>Files in your bucket</h3>
      <ul style={{ listStyle: 'none', padding: 0 }}>
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
            }}
          >
            {file.title}
          </li>
        ))}
      </ul>

      {selectedFile && (
        <p style={{ marginTop: '1rem', fontStyle: 'italic' }}>
          "{selectedFile.title}" will be converted to audio.
        </p>
      )}
    </div>
  );

}

export default App
