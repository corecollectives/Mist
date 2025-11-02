
import { create } from 'zustand';
import { subscribeWithSelector } from 'zustand/middleware';
import { projectsApi } from '../../api';
import type { Project, ProjectCreateInput, ProjectUpdateInput } from '../../types';
import { ERROR_MESSAGES } from '../../constants';

interface ProjectState {
  projects: Project[];
  currentProject: Project | null;
  
  isLoading: boolean;
  error: string | null;
  hasFetched: boolean;
  
  fetchProjects: (force?: boolean) => Promise<void>;
  fetchProjectById: (projectId: string | number) => Promise<void>;
  createProject: (projectData: ProjectCreateInput) => Promise<boolean>;
  updateProject: (projectId: number, projectData: ProjectUpdateInput) => Promise<boolean>;
  deleteProject: (projectId: number) => Promise<boolean>;
  clearError: () => void;
  clearCache: () => void;
  setCurrentProject: (project: Project | null) => void;
  
  getProjectById: (projectId: number) => Project | undefined;
  getProjectsByOwner: (ownerId: number) => Project[];
  getProjectsByTag: (tag: string) => Project[];
}

export const useProjectStore = create<ProjectState>()(
  subscribeWithSelector((set, get) => ({
    projects: [],
    currentProject: null,
    isLoading: false,
    error: null,
    hasFetched: false,

    fetchProjects: async (force = false) => {
      const state = get();
      
      if (state.hasFetched && state.projects.length > 0 && !force) {
        return;
      }

      set({ isLoading: true, error: null });

      try {
        const response = await projectsApi.getAll();
        
        set({
          projects: response.data,
          isLoading: false,
          hasFetched: true,
          error: null,
        });
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : ERROR_MESSAGES.GENERIC_ERROR;
        set({
          isLoading: false,
          error: errorMessage,
        });
        throw error;
      }
    },

    fetchProjectById: async (projectId) => {
      const existingProject = get().projects.find(p => p.id.toString() === projectId.toString());
      if (existingProject) {
        set({ currentProject: existingProject });
        return;
      }

      set({ isLoading: true, error: null });

      try {
        const response = await projectsApi.getById(projectId);
        
        set({
          currentProject: response.data,
          isLoading: false,
          error: null,
        });
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : ERROR_MESSAGES.GENERIC_ERROR;
        set({
          isLoading: false,
          error: errorMessage,
        });
        throw error;
      }
    },

    createProject: async (projectData) => {
      set({ error: null });

      try {
        const response = await projectsApi.create(projectData);
        
        set(state => ({
          projects: [...state.projects, response.data],
          error: null,
        }));

        return true;
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : ERROR_MESSAGES.GENERIC_ERROR;
        set({ error: errorMessage });
        return false;
      }
    },

    updateProject: async (projectId, projectData) => {
      set({ error: null });

      try {
        const response = await projectsApi.update(projectId, projectData);
        
        set(state => ({
          projects: state.projects.map(project =>
            project.id === projectId ? response.data : project
          ),
          currentProject: state.currentProject?.id === projectId ? response.data : state.currentProject,
          error: null,
        }));

        return true;
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : ERROR_MESSAGES.GENERIC_ERROR;
        set({ error: errorMessage });
        return false;
      }
    },

    deleteProject: async (projectId) => {
      set({ error: null });

      try {
        await projectsApi.delete(projectId);
        
        set(state => ({
          projects: state.projects.filter(project => project.id !== projectId),
          currentProject: state.currentProject?.id === projectId ? null : state.currentProject,
          error: null,
        }));

        return true;
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : ERROR_MESSAGES.GENERIC_ERROR;
        set({ error: errorMessage });
        return false;
      }
    },

    clearError: () => set({ error: null }),
    
    clearCache: () => set({ 
      projects: [], 
      currentProject: null, 
      hasFetched: false, 
      error: null 
    }),

    setCurrentProject: (project) => set({ currentProject: project }),

    getProjectById: (projectId) => {
      return get().projects.find(project => project.id === projectId);
    },

    getProjectsByOwner: (ownerId) => {
      return get().projects.filter(project => project.ownerId === ownerId);
    },

    getProjectsByTag: (tag) => {
      return get().projects.filter(project => 
        project.tags?.some(projectTag => 
          projectTag.toLowerCase() === tag.toLowerCase()
        )
      );
    },
  }))
);
