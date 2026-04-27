import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import ProfileView from "../views/ProfileView.vue"
import SettingsPage from "../views/SettingsPage.vue"


const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/',redirect: '/login'},
		{path: '/login', component: LoginView},
		{path: "/users/:username",component: ProfileView},
		{path: "/settings",component: SettingsPage},
		{path: "/:catchAll(.*)",component: LoginView}

	]
})

export default router
