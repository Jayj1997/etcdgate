<!--
 * @Author       : jayj
 * @Date         : 2021-12-13 15:08:25
 * @Description  :
 * @LastEditors  : jayj
 * @LastEditTime : 2021-12-13 15:36:25
-->
<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script>
export default {
  name: 'App',
  data: () => ({
    drawer: null,
    links: [
      ['mdi-cog', 'Conf'],
      ['mdi-account-details', 'User'],
      ['mdi-account-alert', 'Role'],
      ['mdi-history', 'History']
    ]
  })
}
</script>

<style lang="scss">
</style>
