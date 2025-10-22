import { useEffect, useRef, useState } from 'react'
import './App.css'

function App() {
  const wsRef = useRef<WebSocket | null>(null)
  const [msg, setMsg] = useState<any>(null)
  
  useEffect(() => {
    console.log("Opening WebSocket connection");
    wsRef.current = new WebSocket('/api/ws');
    
    wsRef.current.onopen = () => {
      console.log("WebSocket connection opened");
      wsRef.current?.send("hello from client");
    }
    
    wsRef.current.onmessage = (event:MessageEvent) => {
      console.log("Received message:", event.data);
      setMsg(event.data);
    }
    
    return () => {
      wsRef.current?.close();
    }
  }, []);
  
  return (
    <>
      <p>this msg is coming from server via ws</p>
      <div>msg: {msg}</div>
    </>
  )
}

export default App
