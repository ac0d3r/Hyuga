import Vue from 'vue'
import Antd from "ant-design-vue";
import axios from 'axios'
import VueAxios from 'vue-axios'
import App from './App.vue'
import VueClipboard from 'vue-clipboard2'
import VueCookies from 'vue-cookies'

import 'ant-design-vue/dist/antd.css';

VueClipboard.config.autoSetContainer = true;
Vue.use(VueClipboard);
Vue.use(VueCookies);
Vue.use(Antd, VueAxios, axios);

Vue.config.productionTip = false
Vue.prototype.COMMON = global

new Vue({
  render: h => h(App),
}).$mount('#app')
