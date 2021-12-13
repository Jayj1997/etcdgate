<!--
 * @Author       : jayj
 * @Date         : 2021-12-13 15:08:25
 * @Description  :
 * @LastEditors  : jayj
 * @LastEditTime : 2021-12-13 15:29:03
-->
<template>
  <v-app>
    <v-navigation-drawer v-model="drawer" permanent app>
      <v-sheet class="pa-4">
        <v-avatar class="mb-4" size="64">
          <img src="../src/assets/logo.png" alt="logo" />
        </v-avatar>
        <div class="text-h5">ETCD — Gate</div>
      </v-sheet>
      <v-divider></v-divider>
      <v-list>
        <v-list-item v-for="[icon, title] in links" :key="icon" link>
          <v-list-item-icon>
            <v-icon class="ml-2">{{ icon }}</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>{{ title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
    <v-main>
      <router-view />
    </v-main>
  </v-app>
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
