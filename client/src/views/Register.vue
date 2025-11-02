<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-50">
    <div class="w-full max-w-md p-8 space-y-6 bg-white rounded-xl shadow-lg">
      <h1 class="text-3xl font-bold text-center text-gray-800">Create an Account</h1>
      <form @submit.prevent="handleRegister" class="space-y-6">
        <div>
          <label for="username" class="block text-sm font-medium text-gray-700">Username</label>
          <input v-model="username" id="username" type="text" required class="w-full px-4 py-3 mt-1 text-gray-700 bg-gray-100 border-transparent rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent">
        </div>
        <div>
          <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
          <input v-model="email" id="email" type="email" required class="w-full px-4 py-3 mt-1 text-gray-700 bg-gray-100 border-transparent rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent">
        </div>
        <div>
          <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
          <input v-model="password" id="password" type="password" required class="w-full px-4 py-3 mt-1 text-gray-700 bg-gray-100 border-transparent rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent">
        </div>
        <button type="submit" :disabled="loading" class="w-full px-4 py-3 font-semibold text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 transition-colors">
          Register
        </button>
        <p v-if="error" class="mt-4 text-sm text-red-600">{{ error.message }}</p>
      </form>
      <p class="mt-6 text-sm text-center text-gray-600">
        Already have an account? 
        <router-link to="/login" class="font-medium text-indigo-600 hover:text-indigo-500">Login</router-link>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useRegister } from '@/composables/useAuth';

const username = ref('');
const email = ref('');
const password = ref('');

const router = useRouter();
const { mutate: register, loading, error, onDone } = useRegister();

const handleRegister = () => {
  register({ username: username.value, email: email.value, password: password.value });
};

onDone(() => {
  router.push('/login');
});
</script>