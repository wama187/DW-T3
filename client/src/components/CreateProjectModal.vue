<template>
  <div class="fixed inset-0 z-10 flex items-center justify-center overflow-y-auto bg-gray-900 bg-opacity-50">
    <div class="w-full max-w-md p-8 space-y-6 bg-white rounded-xl shadow-lg">
      <h2 class="text-2xl font-bold text-center text-gray-800">New Project</h2>
      <form @submit.prevent="handleCreate" class="space-y-6">
        <div>
          <label for="projectName" class="block text-sm font-medium text-gray-700">Project Name</label>
          <input v-model="name" id="projectName" type="text" required class="w-full px-4 py-3 mt-1 text-gray-700 bg-gray-100 border-transparent rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent">
        </div>
        <div class="flex justify-end space-x-4">
          <button type="button" @click="$emit('close')" class="px-4 py-2 font-semibold text-gray-700 bg-gray-200 rounded-lg hover:bg-gray-300 transition-colors">Cancel</button>
          <button type="submit" :disabled="loading" class="px-4 py-2 font-semibold text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 disabled:opacity-50 transition-colors">Create</button>
        </div>
        <p v-if="error" class="mt-4 text-sm text-red-600">{{ error.message }}</p>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useCreateProjectMutation } from '@/composables/useProjects';
import { useProjectsStore } from '@/store/projects';

const emit = defineEmits(['close']);
const name = ref('');
const projectsStore = useProjectsStore();

const { mutate: createProject, loading, error, onDone } = useCreateProjectMutation();

const handleCreate = () => {
  createProject({ name: name.value });
};

onDone(result => {
  projectsStore.addProject(result.data.createProject);
  emit('close');
});
</script>
