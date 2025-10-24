import Loading from "./components/Loading";
import { useAuth } from "./context/AuthContext";
import { Navigate, Route, BrowserRouter as Router, Routes } from "react-router-dom"
import { SetupPage } from "./pages/Setup";
import { HomePage } from "./pages/Home";


export default function App() {
  const { setupRequired } = useAuth();
  console.log("Setup required:", setupRequired);

  if (setupRequired === null) {
    return <div className="flex h-screen w-screen items-center justify-center">
      <Loading />
    </div>;

  }
  return (
    <Router>
      <Routes>
        {setupRequired ? (
          <>
            <Route path="/setup" element={<SetupPage />} />
            <Route path="*" element={<Navigate to="/setup" replace />} />
          </>

        ) :
          (
            <>
              <Route path="*" element={<HomePage />} />
              <Route path="/setup" element={<Navigate to="/" replace />} />
            </>
          )
        }
      </Routes>
    </Router>

  )
}

