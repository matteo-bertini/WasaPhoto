import {createRouter, createWebHashHistory} from 'vue-router'
import LoginPage from '../views/LoginPage.vue'
import ProfilePage from "../views/ProfilePage.vue"
import ProfileNotFoundPage from "../views/ProfileNotFoundPage.vue"
import SettingsPage from "../views/SettingsPage.vue"

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/',redirect: '/login'},
		{path: '/login', component: LoginPage},
		{path: "/users/:Username",component: ProfilePage},
		{path: "/users/:Username/settings",component: SettingsPage},
		{path: "/profilenotfound",component: ProfileNotFoundPage}

	]
})

export default router
