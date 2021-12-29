<!--
 * @Author       : jayj
 * @Date         : 2021-12-13 15:12:52
 * @Description  :
 * @LastEditors  : jayj
 * @LastEditTime : 2021-12-15 17:04:26
-->
<template>
  <v-row no-gutters>
    <v-col cols="3">
      <conf-tree :conf="test_directory.directory" />
    </v-col>
    <v-col cols="9">
      <conf-formatter :code="code" />
    </v-col>
  </v-row>
</template>

<script>
import ConfTree from '../components/ConfTree.vue'
import ConfFormatter from '../components/ConfFormatter.vue'
export default {
  name: 'Conf',
  data () {
    return {
      test_directory: {
        directory: {
          '/test1': {
            perm_type: 2,
            create_revision: 3,
            mod_revision: 3,
            remaining_lease: 0
          },
          '/test1/test11': {
            perm_type: 2,
            create_revision: 6,
            mod_revision: 6,
            remaining_lease: 0
          },
          '/test1/test11/test111': {
            perm_type: 2,
            create_revision: 10,
            mod_revision: 10,
            remaining_lease: 0
          },
          '/test1/test11/test111/test1111/test1111112312412412412toooooooooolongdatatest': {
            perm_type: 2,
            create_revision: 7,
            mod_revision: 7,
            remaining_lease: 0
          },
          '/test1/test12': {
            perm_type: 2,
            create_revision: 7,
            mod_revision: 7,
            remaining_lease: 0
          },
          '/test1/test12/test121': {
            perm_type: 2,
            create_revision: 7,
            mod_revision: 7,
            remaining_lease: 0
          },
          '/test1/test12/test122': {
            perm_type: 2,
            create_revision: 7,
            mod_revision: 7,
            remaining_lease: 0
          },
          '/test1/test12/test123': {
            perm_type: 2,
            create_revision: 7,
            mod_revision: 7,
            remaining_lease: 0
          },
          '/test1/test124': {
            perm_type: 2,
            create_revision: 7,
            mod_revision: 7,
            remaining_lease: 0
          },
          '/test1/test124/test': {
            perm_type: 2,
            create_revision: 7,
            mod_revision: 7,
            remaining_lease: 0
          },
          '/test1/test13': {
            perm_type: 2,
            create_revision: 7,
            mod_revision: 7,
            remaining_lease: 0
          },
          '/test2': {
            perm_type: 2,
            create_revision: 4,
            mod_revision: 4,
            remaining_lease: 0
          },
          '/test2/test21': {
            perm_type: 2,
            create_revision: 8,
            mod_revision: 8,
            remaining_lease: 0
          },
          '/test2/test22': {
            perm_type: 2,
            create_revision: 9,
            mod_revision: 9,
            remaining_lease: 0
          },
          '/test3': {
            perm_type: 2,
            create_revision: 5,
            mod_revision: 5,
            remaining_lease: 0
          },
          '/ABC': {
            perm_type: 2,
            create_revision: 5,
            mod_revision: 5,
            remaining_lease: 0
          },
          '/abc': {
            perm_type: 2,
            create_revision: 5,
            mod_revision: 5,
            remaining_lease: 0
          }
        },
        is_more: false,
        total: 8
      },
      code: 'test'
    }
  },
  components: {
    ConfTree,
    ConfFormatter
  },
  mounted () {
  },
  methods: {
  }
}
</script>

<style scoped>
</style>
