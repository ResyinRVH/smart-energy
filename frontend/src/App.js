import React, { useState } from "react";
import axios from "axios";



function App() {
  const [file, setFile] = useState(null);
  const [query, setQuery] = useState("");
  const [chatQuery, setChatQuery] = useState("");
  const [response, setResponse] = useState("");

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleUpload = async () => {
    if (!file) {
      console.error("No file selected");
      return;
    }
    
    const formData = new FormData();
    formData.append("file", file);
    formData.append("query", query);

    try {
      const res = await axios.post('http://localhost:8080/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      console.log("Response:", res.data);
      setResponse(res.data.answer); // Assuming the response has an 'answer' field
    } catch (error) {
      console.error('Error uploading file:', error);
    }
  };

  const handleChat = async () => {
    try {
      const res = await axios.post("http://localhost:8080/chat", { query: chatQuery, });
      setResponse(res.data.answer);
    } catch (error) {
      console.error("Error querying chat:", error);
    }
  };

  

  return (
  
    <div 
    style={{ 
      maxWidth: "600px", 
      margin: "0 auto", 
      padding: "20px", 
      textAlign: "center", 
      fontFamily: "Arial, sans-serif" 
      }}
    >
      <h1 style={{ color: "#0ff", marginBottom: "20px", textShadow: "0 0 10px #0ff" }}>
        Data Analysis Chatbot
      </h1>
      <div style={{ marginBottom: "20px" }}>
        <input 
        type="file" 
        onChange={handleFileChange} 
        style={{ padding: "10px", 
        marginRight: "10px", 
        border: "1px solid #ccc", 
        borderRadius: "4px" }} />
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Custom query for file analysis..."
          style={{
            padding: "10px",
            marginRight: "10px",
            border: "1px solid #0ff",
            borderRadius: "4px",
            width: "calc(100% - 170px)",
            backgroundColor: "#fff",
            color: "#000000",
          }}
        />
        <button onClick={handleUpload} style={{ padding: "10px 20px", backgroundColor: "#007bff", color: "white", border: "none", borderRadius: "4px", cursor: "pointer" }}>
          Upload and Analyze
        </button>
      </div>
      <div style={{ marginBottom: "20px" }}>
        <input
          type="text"
          value={chatQuery}
          onChange={(e) => setChatQuery(e.target.value)}
          placeholder="Ask a question..."
          style={{ padding: "10px", marginRight: "10px", border: "1px solid #ccc", borderRadius: "4px", width: "calc(100% - 140px)" }}
        />
        <button onClick={handleChat} style={{ padding: "10px 20px", backgroundColor: "#007bff", color: "white", border: "none", borderRadius: "4px", cursor: "pointer" }}>
          Chat
        </button>
      </div>
      <div ox style={{ marginTop: "20px", padding: "10px", border: "1px solid #ccc", borderRadius: "4px", backgroundColor: "#f9f9f9" }}>
        <h2>Response</h2>
        <p>{response}</p>
      </div>
    </div>
  
  );
}

export default App;