import "babel-core/register";
import "babel-polyfill";

import Vue from "vue";
import App from "./main.vue";

import RestApi from "./rest-api";
window.api = new RestApi("/");

window.app = new Vue({
  el: "main",
  render: h => h(App)
});
