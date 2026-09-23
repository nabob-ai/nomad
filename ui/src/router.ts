import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";

const routes: RouteRecordRaw[] = [
  {
    path: "/",
    name: "Home",
    component: () => import("./views/Home.vue"),
  },
  {
    path: "/aidea",
    name: "AI Assistant",
    component: () => import("./views/AIdea.vue"),
  },
  {
    path: "/passport",
    name: "Passport",
    component: () => import("./views/Passport.vue"),
  },
  {
    path: "/agents",
    name: "Agents",
    component: () => import("./views/AgentManagement.vue"),
  },
  {
    path: "/workflows",
    name: "Workflows",
    component: () => import("./views/WorkflowManagement.vue"),
  },
  {
    path: "/rooms",
    name: "Rooms",
    component: () => import("./views/RoomManagement.vue"),
  },
  {
    path: "/landing",
    name: "Landing",
    component: () => import("./views/LandingPage.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
