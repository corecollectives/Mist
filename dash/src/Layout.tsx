import { Outlet } from 'react-router-dom';
import { Sidebar } from './components/Sidebar';

export const Layout = () => {
  return (
    <div className="flex bg-[#0D1117] min-h-screen">
      <Sidebar />
      <main className="ml-64 w-full p-6">
        <Outlet />
      </main>
    </div>
  );
};
