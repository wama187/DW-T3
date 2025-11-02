<template>
  <div class="fixed inset-0 z-10 flex items-center justify-center overflow-y-auto bg-gray-900 bg-opacity-50">
    <div class="w-full max-w-md p-8 space-y-6 bg-white rounded-xl shadow-lg">
      <h2 class="text-2xl font-bold text-center text-gray-800">Delete Project</h2>
      <p class="text-center text-gray-600">Are you sure you want to delete the project "{{ project.name }}"?</p>
      <div class="flex justify-center space-x-4">
        <button @click="$emit('close')" class="px-4 py-2 font-semibold text-gray-700 bg-gray-200 rounded-lg hover:bg-gray-300 transition-colors">Cancel</button>
        <button @click="handleDelete" :disabled="loading" class="px-4 py-2 font-semibold text-white bg-red-600 rounded-lg hover:bg-red-700 disabled:opacity-50 transition-colors">Delete</button>
      </div>
      <p v-if="error" class="mt-4 text-sm text-red-600">{{ error.message }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';
import { useDeleteProjectMutation } from '@/composables/useProjects';
import { useProjectsStore } from '@/store/projects';
import type { Project } from '@/types';

const props = defineProps<{ project: Project }>();
const emit = defineEmits(['close']);

const projectsStore = useProjectsStore();
const { mutate: deleteProject, loading, error, onDone } = useDeleteProjectMutation();

const handleDelete = () => {
  deleteProject({ id: props.project.id });
};

onDone(() => {
  projectsStore.removeProject(props.project.id);
  emit('close');
});
</script>
