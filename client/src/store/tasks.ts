import { defineStore } from 'pinia';
import type { Task } from '@/types';

export const useTasksStore = defineStore('tasks', {
  state: () => ({
    tasks: [] as Task[],
  }),
  actions: {
    setTasks(tasks: Task[]) {
      this.tasks = [...tasks];
    },
    addTask(task: Task) {
      this.tasks.push(task);
    },
    updateTask(updatedTask: Task) {
      const index = this.tasks.findIndex(t => t.id === updatedTask.id);
      if (index !== -1) {
        this.tasks[index] = updatedTask;
      }
    },
    removeTask(id: string) {
      this.tasks = this.tasks.filter(t => t.id !== id);
    },
  },
});
