import { useState } from "react";
import { toast } from "react-toastify";
import { useAuth } from "../context/AuthContext";


export const SetupPage: React.FC = () => {

  const [formData, setFormData] = useState({
    email: '',
    username: '',
    password: ''
  });
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const { setSetupRequired, setUser
  } = useAuth();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { id, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [id]: value
    }));
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);
    try {
      const response = await fetch('/api/auth/signup', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
        credentials: 'include'
      });
      const data = await response.json();
      if (!data.success) {
        toast.error(data.error || 'Setup failed');
        setError(data.error || 'Setup failed');
      } else {
        setSetupRequired(false);
        setUser({ ...data.data, isAdmin: data.data.role === "owner" || data.data.role === "admin" });
        toast.success('Admin account created successfully!');
      }

    } catch (error) {
      console.error('Setup error:', error);
    }
    setIsLoading(false);
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-[#0D1117]">
      <h1 className="text-2xl font-bold text-[#C9D1D9] mb-4">Setup Admin Account</h1>
      {error && (
        <div className="mb-4 p-3 bg-[#F85149] bg-opacity-20 border border-[#F85149] rounded-md text-white">
          {error}
        </div>
      )}
      <form
        onSubmit={handleSubmit}
        className="w-full max-w-sm bg-[#161B22] p-6 rounded-lg border border-[#30363D]"
      >
        <div className="mb-4">
          <label
            className="block text-[#C9D1D9] mb-2"
            htmlFor="email"
          >
            Email
          </label>
          <input
            type="email"
            id="email"
            className="w-full px-3 py-2 bg-[#0D1117] border border-[#30363D] rounded-md 
                                 text-[#C9D1D9] placeholder-gray-500
                                 focus:outline-none focus:ring-2 focus:ring-[#1F6FEB] focus:border-transparent"
            placeholder="Enter your email"
            onChange={handleChange}
            required
          />
        </div>
        <div className="mb-4">
          <label
            className="block text-[#C9D1D9] mb-2"
            htmlFor="username"
          >
            Username
          </label>
          <input
            type="text"
            id="username"
            className="w-full px-3 py-2 bg-[#0D1117] border border-[#30363D] rounded-md 
                                 text-[#C9D1D9] placeholder-gray-500
                                 focus:outline-none focus:ring-2 focus:ring-[#1F6FEB] focus:border-transparent"
            placeholder="Choose a username"
            onChange={handleChange}
            required
          />
        </div>
        <div className="mb-6">
          <label
            className="block text-[#C9D1D9] mb-2"
            htmlFor="password"
          >
            Password
          </label>
          <input
            type="password"
            id="password"
            className="w-full px-3 py-2 bg-[#0D1117] border border-[#30363D] rounded-md 
                                 text-[#C9D1D9] placeholder-gray-500
                                 focus:outline-none focus:ring-2 focus:ring-[#1F6FEB] focus:border-transparent"
            placeholder="Create a strong password"
            onChange={handleChange}
            required
          />
        </div>
        <button
          type="submit"
          className="w-full bg-[#1F6FEB] hover:bg-[#1958BD] text-[#C9D1D9] 
                             font-semibold py-2 px-4 rounded-md
                             focus:outline-none focus:ring-2 focus:ring-[#1F6FEB] focus:ring-offset-2 
                             focus:ring-offset-[#161B22] transition duration-200"
        >
          {isLoading ? 'Setting up...' : 'Create Admin Account'}
        </button>
      </form>
      <p className="mt-4 text-sm text-[#C9D1D9]">
        This will create your admin account
      </p>
    </div>
  );
};
