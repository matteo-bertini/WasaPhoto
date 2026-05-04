import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import ProfileView from "../views/ProfileView.vue"
import SettingsView from "../views/SettingsView.vue"
import HomeView from '../views/HomeView.vue'


const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/',redirect: '/login'},
		{path: '/login', component: LoginView},
		{path: "/users/:username",component: ProfileView},
		{path: "/settings",component: SettingsView},
		{path: "/home",component: HomeView},
		{path: "/:catchAll(.*)",component: LoginView}

	]
})

export default router
