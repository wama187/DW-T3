<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-50">
    <div class="w-full max-w-md p-8 space-y-6 bg-white rounded-xl shadow-lg">
      <h1 class="text-3xl font-bold text-center text-gray-800">Login</h1>
      <form @submit.prevent="handleLogin" class="space-y-6">
        <div>
          <label for="username" class="block text-sm font-medium text-gray-700">Username</label>
          <input v-model="username" id="username" type="text" required class="w-full px-4 py-3 mt-1 text-gray-700 bg-gray-100 border-transparent rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent">
        </div>
        <div>
          <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
          <input v-model="password" id="password" type="password" required class="w-full px-4 py-3 mt-1 text-gray-700 bg-gray-100 border-transparent rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent">
        </div>
        <button type="submit" :disabled="loading" class="w-full px-4 py-3 font-semibold text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 transition-colors">
          Login
        </button>
        <p v-if="error" class="mt-4 text-sm text-red-600">{{ error.message }}</p>
      </form>
      <p class="mt-6 text-sm text-center text-gray-600">
        Don't have an account? 
        <router-link to="/register" class="font-medium text-indigo-600 hover:text-indigo-500">Register</router-link>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useLogin } from '@/composables/useAuth';
import { useAuthStore } from '@/store/auth';

const username = ref('');
const password = ref('');

const router = useRouter();
const authStore = useAuthStore();
const { mutate: login, loading, error, onDone } = useLogin();

const handleLogin = () => {
  login({ username: username.value, password: password.value });
};

onDone(result => {
  const { token, user } = result.data.login;
  authStore.setAuth(token, user);
  router.push('/');
});
</script>