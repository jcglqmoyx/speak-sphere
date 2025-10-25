import {createRouter, createWebHashHistory} from 'vue-router'
import wordRoutes from "@/router/word";
import authRoutes from "@/router/auth";
import settingRoutes from "@/router/setting";
import QueryView from "@/views/query/QueryView.vue";

const routes = [
    {
        path: '/query',
        name: 'query',
        component: QueryView,
    },
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
        ...routes,
    ]
})

export default router
