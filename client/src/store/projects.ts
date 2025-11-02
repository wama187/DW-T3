import { defineStore } from 'pinia';
import type { Project } from '@/types';

export const useProjectsStore = defineStore('projects', {
  state: () => ({
    projects: [] as Project[],
  }),
  actions: {
    setProjects(projects: Project[]) {
      this.projects = [...projects];
    },
    addProject(project: Project) {
      this.projects.push(project);
    },
    updateProject(updatedProject: Project) {
      const index = this.projects.findIndex(p => p.id === updatedProject.id);
      if (index !== -1) {
        this.projects[index] = updatedProject;
      }
    },
    removeProject(id: string) {
      this.projects = this.projects.filter(p => p.id !== id);
    },
  },
});
