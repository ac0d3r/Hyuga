import Vue from "vue";
import VueRouter from "vue-router";
Vue.use(VueRouter);

// 按需导入
// const Home = () => import("@/view/Home");
const AppViewer = () => import("@/view/AppViewer");
const AppDnsLog = () => import("@/view/apps/DnsLog");

const routes = [
  // {
  //   path: "/",
  //   component: Home
  // },
  {
    path: "/",
    name: "AppViewer",
    component: AppViewer,
    children: [
      {
        path: "/",
        component: AppDnsLog
      }
    ]
  }
];

const router = new VueRouter({
  routes
});

export default router;
