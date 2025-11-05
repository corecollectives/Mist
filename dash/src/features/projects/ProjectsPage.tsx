
import React, { useEffect, useState, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'sonner';

import { useAuth } from '../../context/AuthContext';

import { Button } from '../../components/ui/button';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../../components/ui/card';
import { Badge } from '../../components/ui/badge';
import { Separator } from '../../components/ui/separator';
import { FormModal } from '../../components/FormModal';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '../../components/ui/dropdown-menu';
import { MoreVertical, Plus, Search } from 'lucide-react';
import { Input } from '../../components/ui/input';
import { FullScreenLoading } from '../../shared/components';

import { getProjectOwnerDisplay, formatProjectDate, canEditProject, canDeleteProject, validateProjectData, filterProjects, sortProjects } from '../projects/utils';
import type { Project, ProjectCreateInput } from '../../types';
import { ROUTES, SUCCESS_MESSAGES, ERROR_MESSAGES } from '../../constants';

const ProjectCard: React.FC<{
  project: Project;
  onEdit: (project: Project) => void;
  onDelete: (projectId: number) => void;
  canEdit: boolean;
  canDelete: boolean;
}> = ({ project, onEdit, onDelete, canEdit, canDelete }) => {
  const navigate = useNavigate();
  const ownerDisplay = getProjectOwnerDisplay(project);

  return (
    <Card
      className="relative transition-colors hover:border-primary cursor-pointer group"
      onClick={() => navigate(`${ROUTES.PROJECTS}/${project.id}`)}
    >
      <CardHeader className="flex flex-row items-start justify-between space-y-0">
        <div className="min-w-0 flex-1">
          <CardTitle className="truncate">{project.name}</CardTitle>
          <CardDescription className="truncate">
            {project.description}
          </CardDescription>
        </div>

        {(canEdit || canDelete) && (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button 
                variant="ghost" 
                size="icon" 
                className="h-8 w-8 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
                onClick={(e) => e.stopPropagation()}
              >
                <MoreVertical className="w-4 h-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              {canEdit && (
                <DropdownMenuItem
                  onClick={(e) => {
                    e.stopPropagation();
                    onEdit(project);
                  }}
                >
                  Edit
                </DropdownMenuItem>
              )}
              {canDelete && (
                <DropdownMenuItem
                  className="text-destructive"
                  onClick={(e) => {
                    e.stopPropagation();
                    onDelete(project.id);
                  }}
                >
                  Delete
                </DropdownMenuItem>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        )}
      </CardHeader>

      <CardContent>
        {/* Tags */}
        {project.tags && project.tags.length > 0 && (
          <div className="flex flex-wrap gap-2 mb-4">
            {project.tags.map((tag) => (
              <Badge
                key={tag}
                variant="secondary"
                className="capitalize text-xs"
              >
                {tag}
              </Badge>
            ))}
          </div>
        )}

        <Separator className="my-4" />

        {/* Owner and date info */}
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between text-xs text-muted-foreground gap-2">
          <div className="flex items-center gap-2">
            <div className="w-6 h-6 rounded-full bg-muted flex items-center justify-center">
              <span className="text-xs font-medium text-foreground">
                {ownerDisplay.initials}
              </span>
            </div>
            <span className="truncate">{ownerDisplay.name}</span>
          </div>
          <span className="shrink-0">
            {formatProjectDate(project, 'created')}
          </span>
        </div>
      </CardContent>
    </Card>
  );
};

const EmptyState: React.FC<{ onCreateProject: () => void; canCreate: boolean }> = ({ 
  onCreateProject, 
  canCreate 
}) => (
  <div className="flex flex-col items-center justify-center py-12">
    <div className="text-center">
      <h3 className="text-lg font-semibold mb-2">No projects yet</h3>
      <p className="text-muted-foreground mb-4">
        Get started by creating your first project
      </p>
      {canCreate && (
        <Button onClick={onCreateProject}>
          <Plus className="w-4 h-4 mr-2" />
          Create first project
        </Button>
      )}
    </div>
  </div>
);

export const ProjectsPage: React.FC = () => {
  const { user } = useAuth();
  
  const [projects, setProjects] = useState<Project[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingProject, setEditingProject] = useState<Project | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [sortBy] = useState<'name' | 'created' | 'updated' | 'owner'>('created');
  const [sortOrder] = useState<'asc' | 'desc'>('desc');

  const canCreateProjects = user?.isAdmin || user?.role === 'owner';
  
  const filteredAndSortedProjects = useMemo(() => {
    const filtered = filterProjects(projects, searchTerm);
    return sortProjects(filtered, sortBy, sortOrder);
  }, [projects, searchTerm, sortBy, sortOrder]);

  const fetchProjects = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await fetch("/api/projects/getAll", {
        credentials: 'include'
      });
      const data = await response.json();
      if (!data.success) throw new Error(data.error || "Failed to fetch projects");
      setProjects(data.data);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to fetch projects";
      setError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchProjects();
    
    document.body.style.overflow = 'hidden';
    return () => {
      document.body.style.overflow = '';
    };
  }, []);

  const handleCreateOrUpdateProject = async (projectData: ProjectCreateInput) => {
    try {
      const validation = validateProjectData(projectData);
      if (!validation.isValid) {
        const firstError = Object.values(validation.errors)[0];
        toast.error(firstError);
        return;
      }

      const isEditing = !!editingProject;
      const url = isEditing ? `/api/projects/update?id=${editingProject.id}` : "/api/projects/create";
      const method = isEditing ? "PUT" : "POST";

      const response = await fetch(url, {
        method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(projectData),
        credentials: 'include'
      });

      const data = await response.json();
      if (!data.success) throw new Error(data.error || "Failed to save project");

      const message = isEditing ? SUCCESS_MESSAGES.PROJECT_UPDATED : SUCCESS_MESSAGES.PROJECT_CREATED;
      toast.success(message);
      setIsModalOpen(false);
      setEditingProject(null);
      fetchProjects(); // Refresh the projects list
    } catch (error) {
      console.error('Error saving project:', error);
      toast.error(error instanceof Error ? error.message : ERROR_MESSAGES.GENERIC_ERROR);
    }
  };

  const handleDeleteProject = async (projectId: number) => {
    if (!confirm('Are you sure you want to delete this project?')) return;

    try {
      const response = await fetch(`/api/projects/delete?id=${projectId}`, {
        method: "DELETE",
        credentials: 'include'
      });
      const data = await response.json();
      if (!data.success) throw new Error(data.error || "Failed to delete project");

      toast.success(SUCCESS_MESSAGES.PROJECT_DELETED);
      fetchProjects(); // Refresh the projects list
    } catch (error) {
      console.error('Error deleting project:', error);
      toast.error(error instanceof Error ? error.message : ERROR_MESSAGES.GENERIC_ERROR);
    }
  };

  const handleEditProject = (project: Project) => {
    setEditingProject(project);
    setIsModalOpen(true);
  };

  const handleCreateProject = () => {
    setEditingProject(null);
    setIsModalOpen(true);
  };

  if (isLoading && projects.length === 0) {
    return <FullScreenLoading />;
  }

  return (
    <div className="flex flex-col h-screen bg-background overflow-hidden">
      {/* Header */}
      <div className="flex flex-col gap-4 py-6 border-b border-border shrink-0">
        <div className="flex flex-col sm:flex-row justify-between sm:items-center gap-4">
          <div>
            <h1 className="text-2xl font-bold text-foreground">Projects</h1>
            <p className="text-muted-foreground">
              Create and manage your projects
            </p>
          </div>
          {canCreateProjects && (
            <Button onClick={handleCreateProject}>
              <Plus className="w-4 h-4 mr-2" />
              New Project
            </Button>
          )}
        </div>

        {/* Search and filters */}
        {projects.length > 0 && (
          <div className="flex flex-col sm:flex-row gap-4">
            <div className="relative flex-1 max-w-md">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
              <Input
                placeholder="Search projects..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-10"
              />
            </div>
            {/* Add more filters here if needed */}
          </div>
        )}
      </div>

      {/* Content */}
      <div className="flex-1 py-6 overflow-y-auto">
        {error && (
          <Card className="border-destructive text-destructive mb-6">
            <CardContent className="p-4 text-center">
              <p>{error}</p>
              <Button 
                variant="outline" 
                size="sm" 
                className="mt-2" 
                onClick={() => fetchProjects()}
              >
                Try Again
              </Button>
            </CardContent>
          </Card>
        )}

        {filteredAndSortedProjects.length === 0 ? (
          <EmptyState 
            onCreateProject={handleCreateProject} 
            canCreate={canCreateProjects} 
          />
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 pb-16">
            {filteredAndSortedProjects.map((project) => (
              <ProjectCard
                key={project.id}
                project={project}
                onEdit={handleEditProject}
                onDelete={handleDeleteProject}
                canEdit={canEditProject(project, user)}
                canDelete={canDeleteProject(project, user)}
              />
            ))}
          </div>
        )}
      </div>

      {/* Create/Edit Modal */}
      <FormModal
        isOpen={isModalOpen}
        onClose={() => {
          setIsModalOpen(false);
          setEditingProject(null);
        }}
        onSubmit={handleCreateOrUpdateProject}
        title={editingProject ? 'Edit Project' : 'Create New Project'}
        fields={[
          {
            name: 'name',
            label: 'Project Name',
            type: 'text',
            required: true,
            defaultValue: editingProject?.name || '',
          },
          {
            name: 'description',
            label: 'Description',
            type: 'textarea',
            defaultValue: editingProject?.description || '',
          },
          {
            name: 'tags',
            label: 'Tags',
            type: 'tags',
            defaultValue: editingProject?.tags || [],
          },
        ]}
      />
    </div>
  );
};

export default ProjectsPage;
