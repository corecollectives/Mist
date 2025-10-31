import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { toast } from "react-toastify";
import Loading from "../components/Loading";
import { FiEdit2, FiTrash2 } from 'react-icons/fi';
import { useAuth } from "../context/AuthContext";
import type { Project } from "../lib/types";
import { EditProjectModal } from "../components/EditProjectModal";


export const ProjectPage = () => {
  const { user } = useAuth();
  const params = useParams();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [project, setProject] = useState<Project | null>(null);
  // const [apps, setApps] = useState<App[]>([]);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);

  const fetchProjectDetails = async () => {
    try {
      const response = await fetch(`/api/projects/getFromId?id=${params.projectId}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include'
      });
      const data = await response.json();

      if (!data.success) {
        toast.error(data.error || 'Failed to fetch project details');
        throw new Error(data.error || 'Failed to fetch project details');
      }

      setProject(data.data);
      // setApps(data.apps);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch project details');
      toast.error('Failed to fetch project details');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!window.confirm('Are you sure you want to delete this project? This action cannot be undone.')) {
      return;
    }

    try {
      const response = await fetch(`/api/projects/${project?.id}`, {
        method: 'DELETE',
        credentials: 'include'
      });
      console.log(response);
      const data = await response.json();

      if (!data.success) {
        toast.error(data.error || 'Failed to delete project');
        throw new Error(data.error || 'Failed to delete project');
      }

      toast.success('Project deleted successfully');
      navigate('/projects');
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to delete project');
    }
  };

  const handleEdit = async (projectData: { name: string; description: string; tags: string[] }) => {
    try {
      const response = await fetch(`/api/projects/${project?.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(projectData),
      });

      const data = await response.json();

      if (!data.success) {
        toast.error(data.error || 'Failed to update project');
        throw new Error(data.error || 'Failed to update project');
      }

      toast.success(data.message || 'Project updated successfully');
      fetchProjectDetails();
      setIsEditModalOpen(false);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to update project');
    }
  };

  useEffect(() => {
    fetchProjectDetails();
  }, [params.projectId]);

  if (loading) return <Loading />;

  if (error) {
    return (
      <div className="min-h-screen bg-[#0D1117] p-6">
        <div className="bg-[#F8514933] border border-[#F85149] text-[#F85149] p-4 rounded-lg">
          {error}
        </div>
      </div>
    );
  }

  if (!project) return null;

  return (
    <div className="min-h-screen bg-[#0D1117] p-6">
      <div className="flex justify-between items-start mb-6">
        <div>
          <h1 className="text-[#C9D1D9] text-2xl font-bold">{project.name}</h1>
          <p className="text-[#8B949E] mt-1">{project.description}</p>
          <div className="flex flex-wrap gap-2 mt-3">
            {project?.tags?.map((tag) => (
              <span key={tag} className="px-2 py-1 text-xs rounded-full bg-[#1F6FEB33] text-[#1F6FEB]">
                {tag}
              </span>
            ))}
          </div>
          <div className="flex items-center gap-4 mt-4 text-[#8B949E] text-sm">
            <div className="flex items-center gap-2">
              <div className="w-5 h-5 rounded-full bg-[#30363D] flex items-center justify-center">
                <span className="text-[#C9D1D9] text-xs">
                  {project.owner?.username[0].toUpperCase()}
                </span>
              </div>
              <span>Created by {project.owner?.username}</span>
            </div>
            <span>·</span>
            <span>Updated {new Date(project.updatedAt || "").toLocaleDateString()}</span>
          </div>
          <div className="mt-6">
            <h2 className="text-[#C9D1D9] text-lg font-semibold mb-3">Members</h2>
            <div className="flex flex-wrap gap-3">
              {project.projectMembers?.map((member) => (
                <div key={member.id} className="flex items-center gap-2 px-3 py-1 bg-[#21262D] rounded-lg">
                  <div className="w-6 h-6 rounded-full bg-[#30363D] flex items-center justify-center text-[#C9D1D9] text-xs">
                    {member.username[0].toUpperCase()}
                  </div>
                  <span className="text-[#C9D1D9] text-sm">{member.username}</span>
                </div>
              ))}

              {/* Add Member Button */}
              {user?.isAdmin && (
                <button
                  // onClick={() => setIsAddMemberModalOpen(true)}
                  className="flex items-center gap-2 px-3 py-1 bg-[#1F6FEB] text-white rounded-lg hover:bg-[#1A73E8] transition-colors text-sm"
                >
                  + Add Member
                </button>
              )}
            </div>
          </div>

        </div>
        {user?.isAdmin && (
          <div className="flex items-center gap-2">
            <button
              onClick={() => setIsEditModalOpen(true)}
              className="px-3 py-2 bg-[#21262D] text-[#C9D1D9] rounded-lg hover:bg-[#30363D] transition-colors inline-flex items-center gap-2"
            >
              <FiEdit2 className="w-4 h-4" />
              Edit
            </button>
            <button
              onClick={() => handleDelete()}
              className="px-3 py-2 bg-[#21262D] text-[#F85149] rounded-lg hover:bg-[#30363D] transition-colors inline-flex items-center gap-2"
            >
              <FiTrash2 className="w-4 h-4" />
              Delete
            </button>
            <button className="px-3 py-2 bg-[#1F6FEB] text-white rounded-lg hover:bg-[#1A73E8] transition-colors inline-flex items-center gap-2">
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
              </svg>
              New App
            </button>
          </div>
        )}
      </div>


      {project && (
        <EditProjectModal
          isOpen={isEditModalOpen}
          onClose={() => setIsEditModalOpen(false)}
          onSubmit={handleEdit}
          project={project}
        />
      )}
    </div>
  );
};
