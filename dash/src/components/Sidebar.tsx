import { Link, useLocation } from 'react-router-dom';
import { MdDashboard } from "react-icons/md";
import { VscGithubProject } from "react-icons/vsc";
import { FaDatabase } from "react-icons/fa";
import { FaRocket } from "react-icons/fa";
import { FaScroll } from "react-icons/fa";
import { MdOutlineLogout } from "react-icons/md";
import { FaUsers } from "react-icons/fa";
import { IoMdSettings } from "react-icons/io";
import { useAuth } from '../context/AuthContext';

const routes = [
  { path: '/', label: 'Dashboard', icon: <MdDashboard className="text-2xl" /> },
  { path: '/projects', label: 'Projects', icon: <VscGithubProject className="text-2xl" /> },
  { path: '/databases', label: 'Databases', icon: <FaDatabase className="text-2xl" /> },
  { path: '/deployments', label: 'Deployments', icon: <FaRocket className='text-2xl' /> },
  { path: '/logs', label: 'Logs', icon: <FaScroll className='text-2xl' /> },
  // { path: '/terminal', label: 'Terminal', icon: <IoTerminal className='text-2xl' /> },
  { path: '/users', label: 'Users', icon: <FaUsers className='text-2xl' /> },
  { path: '/settings', label: 'Settings', icon: <IoMdSettings className='text-2xl' /> },
];

export const Sidebar = () => {
  const location = useLocation();
  const { user, logout } = useAuth();
  return (
    <div className="h-screen w-64 bg-[#161B22] border-r border-[#30363D] fixed left-0 top-0">
      <div className="flex flex-col border-b border-[#30363D] items-start justify-center">
        <div className="p-4 flex items-center gap-x-4">
          <img
            src='/mist.png'
            alt='Logo'
            className="w-16 h-16 object-center"
          />
          <div className="flex flex-col text-3xl text-blue-500 font-bold items-center justify-center">
            Mist
          </div>
        </div>
      </div>
      <nav className="p-4">
        {routes.map((route) => (
          <Link
            key={route.path}
            to={route.path}
            className={`flex items-center p-3 mb-2 rounded-md transition-colors
              ${location.pathname === route.path
                ? 'bg-[#1F6FEB] text-[#C9D1D9]'
                : 'text-[#C9D1D9] hover:bg-[#30363D]'
              }`}
          >
            <span className="mr-3">{route.icon}</span>
            {route.label}
          </Link>
        ))}
      </nav>
      <div className="absolute bottom-0 left-0 w-full justify-between p-4 flex gap-x-5 border-t items-center border-[#30363D]">
        <div className="flex items-center gap-x-2">
          <span className="w-8 h-8 rounded-full text-black bg-blue-500 flex items-center justify-center">
            {user?.username[0].toUpperCase() || 'C'}
          </span>
          <span className="text-[#C9D1D9] text-2xl">{user?.username || 'Guest'}</span>
        </div>
        <MdOutlineLogout
          className="text-2xl right-0 text-[#C9D1D9] cursor-pointer hover:text-red-500"
          onClick={logout}
        />
      </div>

    </div>
  );
};
