<template>
  <router-view></router-view>
</template>

<script>
import PageNotFound from "./components/404";

import Vue from "vue";
import VueRouter from "vue-router";
Vue.use(VueRouter);

const router = new VueRouter({
  mode: "history",
  routes: [
    // ...
    { path: "*", component: PageNotFound }
  ]
});

export default {
  router
};
</script>
