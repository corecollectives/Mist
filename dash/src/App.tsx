import './App.css'

function App() {
  const fetchData = async () => {
    const response = await fetch('/api/health', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },

    })
    const data = await response.json()
    console.log(data)
  }
  fetchData()
  return (
    <>
      Hello
    </>
  )
}

export default App
