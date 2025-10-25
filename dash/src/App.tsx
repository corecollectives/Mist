import Loading from "./components/Loading";
import { useAuth } from "./context/AuthContext";
import { Navigate, Route, BrowserRouter as Router, Routes } from "react-router-dom"
import { SetupPage } from "./pages/Setup";
import { HomePage } from "./pages/Home";
import { LoginPage } from "./pages/Login";
import { Layout } from "./Layout";
import { UsersPage } from "./pages/Users";
import { ProjectsPage } from "./pages/Projects";
import { DeploymentsPage } from "./pages/Deployments";
import { DatabasesPage } from "./pages/Databases";
import { LogsPage } from "./pages/Logs";
import { SettingsPage } from "./pages/Settings";

export default function App() {
  const { setupRequired, user } = useAuth();

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
        ) : !user ? (
          <>
            <Route path="/login" element={<LoginPage />} />
            <Route path="*" element={<Navigate to="/login" replace />} />
          </>
        ) : (
          <>
            <Route element={<Layout />} >
              <Route path="/" element={<HomePage />} />
              <Route path="/users" element={<UsersPage />} />
              <Route path="/projects" element={<ProjectsPage />} />
              <Route path="/deployments" element={<DeploymentsPage />} />
              <Route path="/databases" element={<DatabasesPage />} />
              <Route path="/logs" element={<LogsPage />} />
              <Route path="/settings" element={<SettingsPage />} />
            </Route>
            <Route path="*" element={<Navigate to="/" replace />} />
          </>
        )}
      </Routes>
    </Router>
  )
}

