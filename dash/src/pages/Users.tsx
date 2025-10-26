import { useEffect, useState } from "react";
import type { User } from "../lib/types";
import { useAuth } from "../context/AuthContext";
import Loading from "../components/Loading";
import { toast } from "react-toastify";
import { CreateUserModal } from "../components/CreateUserModal";

const getRoleStyles = (role: string) => {
  switch (role) {
    case 'owner':
      return 'bg-[#A371F733] text-[#A371F7]';
    case 'admin':
      return 'bg-[#1F6FEB33] text-[#1F6FEB]';
    default:
      return 'bg-[#30363D] text-[#8B949E]';
  }
};

export const UsersPage = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const { user } = useAuth();

  const fetchUsers = async () => {
    try {
      const response = await fetch('/api/users/getAll');
      const data = await response.json();

      if (!data.success) {
        throw new Error(data.error || 'Failed to fetch users');
      }
      const updatedUsers: User[] = data.data.map((u: User) => {
        u.isAdmin = u.role === 'admin' || u.role === 'owner';
        return u;
      })
      setUsers(updatedUsers);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch users');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const handleCreateUser = async (userData: { username: string; email: string; password: string; role: 'admin' | 'user' }) => {
    try {
      const response = await fetch('/api/users/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData),
      });

      const data = await response.json();

      if (!data.success) {
        toast.error(data.error || 'Failed to create user');
      }
      toast.success(data.message || 'User created successfully');
      fetchUsers();
      setIsModalOpen(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create user');
    }
  };

  if (loading) {
    return <div className="h-full w-full">
      <Loading />
    </div>;
  }
  return (
    <div className="min-h-screen bg-[#0D1117] p-6">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-[#C9D1D9] text-2xl font-bold">Users</h1>
          <p className="text-[#8B949E] mt-1">Manage users and their permissions</p>
        </div>
        <div className="flex items-center gap-4">
          <button
            disabled={!user?.isAdmin}
            className="px-4 py-2 rounded-lg text-white transition-colors
               bg-[#1F6FEB] hover:bg-[#1A73E8]
               disabled:bg-[#555] disabled:text-[#ccc] disabled:cursor-not-allowed"
            onClick={() => setIsModalOpen(true)}
          >
            Add User
          </button>
        </div>
      </div>

      {error ? (
        <div className="bg-[#F8514933] border border-[#F85149] text-[#F85149] p-4 rounded-lg">
          {error}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {users.map((user) => (
            <div key={user.id} className="bg-[#161B22] cursor-pointer border border-[#30363D] rounded-lg p-4 hover:border-[#1F6FEB] transition-colors">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className='flex gap-x-2 items-center'>
                    <div className="w-8 h-8 rounded-full bg-[#30363D] flex items-center justify-center">
                      <span className="text-[#C9D1D9] text-sm">
                        {user.username[0].toUpperCase()}
                      </span>
                    </div>
                    <h3 className="text-[#C9D1D9] font-semibold text-lg">{user.username}</h3>
                  </div>
                  <p className="text-[#8B949E] text-sm mt-1">{user.email}</p>
                </div>
                <span className={`px-2 py-1 text-xs rounded-full ${getRoleStyles(user.role)}`}>
                  {user.role.charAt(0).toUpperCase() + user.role.slice(1)}
                </span>
              </div>

              <div className="mt-4 pt-4 border-t border-[#30363D]">
                <div className="flex items-center gap-2 mb-3">

                  <span className="text-[#8B949E] text-sm font-mono break-all">
                    User Id: {user.id}
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      <CreateUserModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleCreateUser}
      />
    </div>
  );
}
