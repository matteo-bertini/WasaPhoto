import {createRouter, createWebHashHistory} from 'vue-router'
import LoginPage from '../views/LoginPage.vue'
import ProfilePage from "../views/ProfilePage.vue"

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/',redirect: '/login'},
		{path: '/login', component: LoginPage},
		{path: "/:Username",component: ProfilePage}

	]
})

export default router
