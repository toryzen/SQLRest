// src/router/index.js  
import { createRouter, createWebHistory } from 'vue-router';  
  
// 页面组件  
import ProjectList from '../views/ProjectList.vue';  
import DBList from '../views/DBList.vue';  
import AuthKeyList from '../views/AuthKeyList.vue';  
import APIList from '../views/APIList.vue';  
  
const routes = [  
  {  
    path: '/',  
    redirect: '/project-list',  
  },  
  {  
    path: '/project-list',  
    name: 'ProjectList',  
    component: ProjectList,  
  },  
  {  
    path: '/db-list',  
    name: 'DBList',  
    component: DBList,  
  },  
  {  
    path: '/authkey-list',  
    name: 'AuthKeyList',  
    component: AuthKeyList,  
  },  
  {  
    path: '/api-list',  
    name: 'APIList',  
    component: APIList,  
  },  
];  
  
const router = createRouter({  
  history: createWebHistory(process.env.BASE_URL),  
  routes,  
});  
  
export default router;  
