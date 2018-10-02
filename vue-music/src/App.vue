<template>
    <div id="app">
        <m-header></m-header>
        <tab></tab>
        <keep-alive>
            <router-view></router-view>
        </keep-alive>
    </div>
</template>

<script type="text/ecmascript-6">
import MHeader from 'components/m-header/m-header.vue'
import Tab from 'components/tab/tab.vue'
export default {
  components: {
    MHeader,
    Tab
  }
}
</script>

<style lang="css" rel="stylesheet/css">
body {
  margin: 0;
  padding: 0;
  line-height: 1;
  background-color: #E7D5E5;
}
</style>
