<template>
  <div id="app">
    <Viewer />
  </div>
</template>

<script>
import Viewer from "./view/Viewer.vue";

export default {
  name: "App",
  components: {
    Viewer,
  },
};
</script>

<style>
body {
  font-size: 16px;
  -webkit-font-smoothing: antialiased;
}
</style>