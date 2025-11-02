<template>
  <div class="fixed inset-0 z-10 flex items-center justify-center overflow-y-auto bg-gray-900 bg-opacity-50">
    <div class="w-full max-w-md p-8 space-y-6 bg-white rounded-xl shadow-lg">
      <h2 class="text-2xl font-bold text-center text-gray-800">Edit Task</h2>
      <form @submit.prevent="handleUpdate" class="space-y-6">
        <div>
          <label for="taskTitle" class="block text-sm font-medium text-gray-700">Task Title</label>
          <input v-model="title" id="taskTitle" type="text" required class="w-full px-4 py-3 mt-1 text-gray-700 bg-gray-100 border-transparent rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent">
        </div>
        <div>
          <label for="taskStatus" class="block text-sm font-medium text-gray-700">Status</label>
          <select v-model="status" id="taskStatus" class="w-full px-4 py-3 mt-1 text-gray-700 bg-gray-100 border-transparent rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent">
            <option value="PENDING">Pending</option>
            <option value="IN_PROGRESS">In Progress</option>
            <option value="COMPLETED">Completed</option>
          </select>
        </div>
        <div class="flex justify-end space-x-4">
          <button type="button" @click="$emit('close')" class="px-4 py-2 font-semibold text-gray-700 bg-gray-200 rounded-lg hover:bg-gray-300 transition-colors">Cancel</button>
          <button type="submit" :disabled="loading" class="px-4 py-2 font-semibold text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 disabled:opacity-50 transition-colors">Update</button>
        </div>
        <p v-if="error" class="mt-4 text-sm text-red-600">{{ error.message }}</p>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, defineProps, defineEmits, watch } from 'vue';
import { useUpdateTaskMutation } from '@/composables/useTasks';
import { useTasksStore } from '@/store/tasks';
import type { Task } from '@/types';

const props = defineProps<{ task: Task }>();
const emit = defineEmits(['close']);

const title = ref(props.task.title);
const status = ref(props.task.status);
const tasksStore = useTasksStore();

watch(
  () => props.task,
  (task) => {
    title.value = task.title;
    status.value = task.status;
  },
  { immediate: true }
);


const { mutate: updateTask, loading, error, onDone } = useUpdateTaskMutation();

const handleUpdate = () => {
  updateTask({ id: props.task.id, title: title.value, status: status.value });
};

onDone(result => {
  tasksStore.updateTask(result.data.updateTask);
  emit('close');
});
</script>
