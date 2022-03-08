<template>
  <a-layout id="app-viewer">
    <a-layout-header class="container-page-header">
      <div class="container-page-header-logo">
        <img :src="LogoImg" alt="logo" />
      </div>
    </a-layout-header>
    <a-layout-content>
      <div>
        <DnsLog />
      </div>
    </a-layout-content>
    <a-layout-footer :style="{ textAlign: 'center' }">
      <a
        target="_blank"
        href="https://github.com/Buzz2d0/Hyuga"
        style="padding: 0px 5px 0px 5px"
        ><a-icon type="github" style="color: rgb(61, 64, 67)"/>
        Hyuga
        <a-icon type="star" theme="filled" style="color: rgb(226, 197, 98)"
      /></a>

      © 2019-2022 Created by zznQ
      <a-icon type="sync" spin />
    </a-layout-footer>
  </a-layout>
</template>

<script>
import DnsLog from "../components/DnsLog.vue";
export default {
  name: "Viewer",
  components: {
    DnsLog,
  },
  data: () => ({
    LogoImg: require("../assets/logo.png"),
  }),
};
</script>

<style lang="scss" scoped>
.container-page-header {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  background-color: #ffffff;
  &-logo {
    img {
      max-width: 100px;
    }
  }
}
</style>
