import {createRouter, createWebHashHistory} from 'vue-router'
import vocabularyRoutes from "@/router/vocabulary";
import authRoutes from "@/router/auth";
import settingRoutes from "@/router/setting";
import toolRoutes from "@/router/tool";

const routes = [
    {
        path: '/',
        name: 'home',
        redirect: '/vocabulary',
    },
    {
        path: '/404',
        name: '404',
        redirect: '/vocabulary'
    },
    {
        path: '/:pathMatch(.*)*',
        redirect: '/404'
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        ...vocabularyRoutes,
        ...authRoutes,
        ...settingRoutes,
        ...toolRoutes,
        ...routes,
    ]
})

export default router
