import { createApp, provide, h } from 'vue';
import { DefaultApolloClient } from '@vue/apollo-composable';
import App from './App.vue';
import router from './router';
import pinia from './plugins/pinia';
import apolloClient from './plugins/apollo';
import './style.css';

const app = createApp({
  setup() {
    provide(DefaultApolloClient, apolloClient);
  },
  render: () => h(App),
});

app.use(pinia);
app.use(router);

app.mount('#app');