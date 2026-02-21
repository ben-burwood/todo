import { createApp } from 'vue'
import App from './App.vue'
import './style.css';

export const SERVER_URL = `${window.location.protocol}//${window.location.hostname}:${window.location.port}`;
export const EDIT_ENABLED = true;
export const DELETE_ENABLED = true;

createApp(App).mount('#app')
