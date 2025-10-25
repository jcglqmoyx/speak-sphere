import {createRouter, createWebHashHistory} from 'vue-router'
import wordRoutes from "@/router/word";
import authRoutes from "@/router/auth";
import settingRoutes from "@/router/setting";
import toolRoutes from "@/router/tool";

const routes = [
    {
        path: '/',
        name: 'home',
        redirect: '/word',
    },
    {
        path: '/404',
        name: '404',
        redirect: '/word'
    },
    {
        path: '/:pathMatch(.*)*',
        redirect: '/404'
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        ...wordRoutes,
        ...authRoutes,
        ...settingRoutes,
        ...toolRoutes,
        ...routes,
    ]
})

export default router
