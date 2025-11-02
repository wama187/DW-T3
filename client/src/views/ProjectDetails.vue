<template>
  <div class="min-h-screen bg-white">
    <!-- Header con línea de acento -->
    <header class="border-b border-gray-100 relative">
      <div class="absolute bottom-0 left-0 right-0 h-0.5 bg-gradient-to-r from-indigo-500 via-purple-500 to-pink-500"></div>
      <div class="px-4 py-8 mx-auto max-w-7xl sm:px-6 lg:px-8">
        <div class="flex items-center justify-between">
          <div class="space-y-1">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-gradient-to-br from-indigo-100 to-purple-100 rounded-lg flex items-center justify-center">
                <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"></path>
                </svg>
              </div>
              <h1 class="text-4xl font-bold text-gray-900 tracking-tight">{{ project?.name }}</h1>
            </div>
            <p class="text-sm text-gray-500 ml-13">Project workspace</p>
          </div>
          <router-link 
            to="/" 
            class="group px-5 py-2.5 font-medium text-gray-700 bg-white border border-gray-200 rounded-lg hover:bg-gray-50 hover:border-gray-300 transition-all duration-200"
          >
            <span class="flex items-center gap-2">
              <svg class="w-4 h-4 text-gray-500 group-hover:text-indigo-600 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
              </svg>
              Back to Projects
            </span>
          </router-link>
        </div>
      </div>
    </header>

    <main class="py-16 mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <!-- Header section -->
      <div class="flex items-center justify-between mb-12">
        <div>
          <h2 class="text-2xl font-bold text-gray-900">Tasks</h2>
          <p class="mt-1 text-sm text-gray-500">{{ tasks?.length || 0 }} task{{ tasks?.length !== 1 ? 's' : '' }} in this project</p>
        </div>
        <button 
          @click="showCreateModal = true" 
          class="group px-6 py-3 font-semibold text-white bg-gradient-to-r from-indigo-600 to-purple-600 rounded-lg hover:from-indigo-700 hover:to-purple-700 transition-all duration-200 shadow-sm hover:shadow-md"
        >
          <span class="flex items-center gap-2">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
            </svg>
            New Task
          </span>
        </button>
      </div>

      <!-- Loading state -->
      <div v-if="loading" class="flex items-center justify-center py-20">
        <div class="text-center space-y-4">
          <div class="inline-flex items-center justify-center">
            <svg class="animate-spin h-10 w-10 text-indigo-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </div>
          <p class="text-gray-500 font-medium">Loading tasks...</p>
        </div>
      </div>

      <!-- Error state -->
      <div v-if="error" class="p-5 bg-red-50 border border-red-100 rounded-xl">
        <div class="flex items-start gap-3">
          <svg class="w-6 h-6 text-red-500 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"></path>
          </svg>
          <div>
            <h3 class="font-semibold text-red-900">Error loading tasks</h3>
            <p class="text-sm text-red-700 mt-1">{{ error.message }}</p>
          </div>
        </div>
      </div>

      <!-- Tasks grid -->
      <div v-if="tasks && tasks.length" class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        <div 
          v-for="task in tasks" 
          :key="task.id" 
          class="group relative p-6 bg-white border border-gray-200 rounded-xl hover:border-indigo-200 hover:shadow-lg transition-all duration-200"
        >
          <!-- Indicador de color superior -->
          <div class="absolute top-0 left-0 right-0 h-1 rounded-t-xl opacity-0 group-hover:opacity-100 transition-opacity" :class="statusBarColor(task.status)"></div>
          
          <div class="space-y-4">
            <!-- Header con título -->
            <div class="flex items-start justify-between gap-3">
              <h3 class="text-xl font-bold text-gray-900 group-hover:text-indigo-600 transition-colors line-clamp-2 flex-1">
                {{ task.title }}
              </h3>
              <div class="flex-shrink-0">
                <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
                </svg>
              </div>
            </div>

            <!-- Status badge -->
            <div class="flex items-center gap-2">
              <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold" :class="statusBadgeClass(task.status)">
                <span class="w-1.5 h-1.5 rounded-full" :class="statusDotColor(task.status)"></span>
                {{ formatStatus(task.status) }}
              </span>
            </div>

            <!-- Action buttons -->
            <div class="flex items-center gap-2 pt-4 border-t border-gray-100">
              <button 
                @click="openEditModal(task)" 
                class="flex-1 px-3 py-2 text-sm font-medium text-gray-700 bg-gray-50 rounded-lg hover:bg-gray-100 hover:text-gray-900 transition-colors"
              >
                <span class="flex items-center justify-center gap-1.5">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                  </svg>
                  Edit
                </span>
              </button>
              <button 
                @click="openDeleteModal(task)" 
                class="px-3 py-2 text-sm font-medium text-red-600 bg-red-50 rounded-lg hover:bg-red-100 hover:text-red-700 transition-colors"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty state -->
      <div v-else-if="!loading && !error" class="text-center py-20">
        <div class="inline-flex items-center justify-center w-20 h-20 bg-gray-100 rounded-full mb-6">
          <svg class="w-10 h-10 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"></path>
          </svg>
        </div>
        <h3 class="text-2xl font-bold text-gray-900 mb-2">No tasks yet</h3>
        <p class="text-gray-500 mb-8 max-w-md mx-auto">Start organizing your work by creating your first task for this project.</p>
        <button 
          @click="showCreateModal = true"
          class="inline-flex items-center gap-2 px-6 py-3 font-semibold text-white bg-gradient-to-r from-indigo-600 to-purple-600 rounded-lg hover:from-indigo-700 hover:to-purple-700 transition-all duration-200 shadow-sm hover:shadow-md"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
          </svg>
          Create Your First Task
        </button>
      </div>
    </main>

    <!-- Modals -->
    <CreateTaskModal v-if="showCreateModal" :project-id="projectId" @close="showCreateModal = false; refreshTasks();" />
    <EditTaskModal v-if="taskToEdit" :task="taskToEdit" @close="taskToEdit = null; refreshTasks();" />
    <DeleteTaskModal v-if="taskToDelete" :task="taskToDelete" @close="taskToDelete = null; refreshTasks();" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRoute } from 'vue-router';
import { useTasksStore } from '@/store/tasks';
import { useProjectQuery } from '@/composables/useProjects';
import CreateTaskModal from '@/components/CreateTaskModal.vue';
import EditTaskModal from '@/components/EditTaskModal.vue';
import DeleteTaskModal from '@/components/DeleteTaskModal.vue';
import type { Task } from '@/types';

const route = useRoute();
const tasksStore = useTasksStore();
const projectId = route.params.id as string;

const { result, loading, error, onResult, refetch } = useProjectQuery(projectId);

onResult(queryResult => {
  if (queryResult.data && queryResult.data.project) {
    tasksStore.setTasks(queryResult.data.project.tasks);
  }
});

const refreshTasks = () => {
  refetch();
};

const project = computed(() => result.value?.project);
const tasks = computed(() => tasksStore.tasks);

const showCreateModal = ref(false);
const taskToEdit = ref<Task | null>(null);
const taskToDelete = ref<Task | null>(null);

const openEditModal = (task: Task) => {
  taskToEdit.value = task;
};

const openDeleteModal = (task: Task) => {
  taskToDelete.value = task;
};

const formatStatus = (status: string) => {
  return status.replace('_', ' ');
};

const statusBadgeClass = (status: string) => {
  switch (status) {
    case 'PENDING':
      return 'bg-yellow-50 text-yellow-700 border border-yellow-200';
    case 'IN_PROGRESS':
      return 'bg-blue-50 text-blue-700 border border-blue-200';
    case 'COMPLETED':
      return 'bg-green-50 text-green-700 border border-green-200';
    default:
      return 'bg-gray-50 text-gray-700 border border-gray-200';
  }
};

const statusDotColor = (status: string) => {
  switch (status) {
    case 'PENDING':
      return 'bg-yellow-500';
    case 'IN_PROGRESS':
      return 'bg-blue-500';
    case 'COMPLETED':
      return 'bg-green-500';
    default:
      return 'bg-gray-500';
  }
};

const statusBarColor = (status: string) => {
  switch (status) {
    case 'PENDING':
      return 'bg-yellow-400';
    case 'IN_PROGRESS':
      return 'bg-blue-400';
    case 'COMPLETED':
      return 'bg-green-400';
    default:
      return 'bg-gray-400';
  }
};
</script>