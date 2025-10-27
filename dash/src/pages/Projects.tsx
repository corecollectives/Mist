

import { useEffect, useState } from 'react';
import Loading from '../components/Loading';
import { toast } from 'react-toastify';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { CreateProjectModal } from '../components/CreateProjectModal';
import type { Project } from '../lib/types';


export const ProjectsPage = () => {
  const { user } = useAuth();
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const navigate = useNavigate();
  const fetchProjects = async () => {
    try {
      const response = await fetch('/api/projects/getAll');
      const data = await response.json();

      if (!data.success) {
        toast.error(data.error || 'Failed to fetch projects');
        throw new Error(data.error || 'Failed to fetch projects');
      }

      setProjects(data.data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch projects');
      toast.error('Failed to fetch projects');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProjects();
  }, []);

  const handleCreateProject = async (projectData: { name: string; description: string; tags: string[] }) => {
    try {
      const response = await fetch('/api/projects/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(projectData),
      });

      const data = await response.json();

      if (!data.success) {
        toast.error(data.error || 'Failed to create project');
        throw new Error(data.error || 'Failed to create project');
      }

      toast.success(data.message || 'Project created successfully');
      fetchProjects();
      setIsModalOpen(false);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to create project');
    }
  };

  if (loading) return <Loading />;
  return (
    <div className="min-h-screen bg-[#0D1117] p-6">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-[#C9D1D9] text-2xl font-bold">Projects</h1>
          <p className="text-[#8B949E] mt-1">Manage your projects and deployments</p>
        </div>
        {user?.isAdmin && (
          <button
            onClick={() => setIsModalOpen(true)}
            className="px-4 cursor-pointer py-2 bg-[#1F6FEB] text-white rounded-lg hover:bg-[#1A73E8] transition-colors"
          >
            New Project
          </button>
        )}
      </div>

      {error ? (
        <div className="bg-[#F8514933] border border-[#F85149] text-[#F85149] p-4 rounded-lg">
          {error}
        </div>
      ) : projects.length === 0 ? (
        <div className="bg-[#161B22] border border-[#30363D] rounded-lg p-8 text-center">
          <p className="text-[#C9D1D9] text-lg mb-4">No projects yet</p>
          {user?.isAdmin && (
            <button
              onClick={() => setIsModalOpen(true)}
              className="px-4 py-2 cursor-pointer bg-[#1F6FEB] text-white rounded-lg hover:bg-[#1A73E8] transition-colors inline-flex items-center gap-2"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
              </svg>
              Create first project
            </button>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {projects.map((project) => (
            <div onClick={() => navigate(`/projects/${project.id}`)} key={project.id} className="bg-[#161B22] cursor-pointer border border-[#30363D] rounded-lg p-4 hover:border-[#1F6FEB] transition-colors">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <h3 className="text-[#C9D1D9] font-semibold text-lg">{project.name}</h3>
                  <p className="text-[#8B949E] text-sm mt-1">{project.description}</p>
                </div>
              </div>

              <div className="mt-3 flex flex-wrap gap-2">
                {project?.tags?.map((tag) => (
                  <span key={tag} className="px-2 py-1 text-xs rounded-full bg-[#1F6FEB33] text-[#1F6FEB]">
                    {tag}
                  </span>
                ))}
              </div>

              <div className="mt-4 pt-4 border-t border-[#30363D]">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <div className="w-6 h-6 rounded-full bg-[#30363D] flex items-center justify-center">
                      <span className="text-[#C9D1D9] text-xs">
                        {project.owner?.username[0].toUpperCase()}
                      </span>
                    </div>
                    <span className="text-[#8B949E] text-sm">
                      {project.owner?.username}
                    </span>
                  </div>
                  <span className="text-[#8B949E] text-sm">
                    {new Date(project.updatedAt || "").toLocaleDateString()}
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      <CreateProjectModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleCreateProject}
      />
    </div>
  );
};


