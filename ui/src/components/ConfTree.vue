<!--
 * @Author       : jayj
 * @Date         : 2021-12-14 10:19:40
 * @Description  :
 * @LastEditors  : jayj
 * @LastEditTime : 2021-12-15 17:28:04
-->
<template>
  <v-card tile>
    <v-sheet class="pa-4 primary lighten-2">
      <v-text-field
        v-model="search"
        label="Search Keys"
        dark
        flat
        solo-inverted
        hide-details
        clearable
        clear-icon="mdi-close-circle-outline"
      ></v-text-field>
      <v-checkbox
        v-model="caseSensitive"
        dark
        hide-details
        label="Case sensitive search"
      ></v-checkbox>
    </v-sheet>
    <v-card-text class="tree">
      <v-treeview
        :items="generateConfigTree"
        :search="search"
        :filter="searchFilter"
        :open.sync="open"
        :active.sync="choseNode"
        color="warning"
        dense
        activatable
      >
      </v-treeview>
    </v-card-text>
  </v-card>
</template>

<script>

export default {
  name: 'ConfTree',
  props: {
    conf: {
      type: Object
    }
  },
  data () {
    return {
      open: [1, 2],
      search: null,
      caseSensitive: false,
      choseNode: [],
      index_id: 1// for index and search
    }
  },
  computed: {
    generateConfigTree () { // generate tree view from props(response)
      const vm = this
      const tree = [
        {
          id: vm.index_id, // for index and search
          name: '/',
          children: {}
        }
      ]
      vm.index_id += 1

      const sub = tree[0]

      const configs = this.$props.conf

      for (const conf in configs) {
        const splitKey = conf.split('/')

        let cur = sub

        for (let i = 0; i < splitKey.length; i++) {
          if (splitKey[i] === '') continue

          if (!cur.children[splitKey[i]]) {
            cur.children[splitKey[i]] = {
              children: {}
            }
          }

          cur = cur.children[splitKey[i]]
        }
      }

      tree[0].children = this.convertObjToArr(tree[0].children)

      return tree
    },
    searchFilter () { // filter search
      return this.caseSensitive
        ? (item, search, textKey) => item[textKey].indexOf(search) > -1
        : undefined
    }
  },
  methods: {
    // => from
    // {
    //   a: {children: {}},
    //   b: {},
    // }
    // => to
    // [
    //   {name: a, children: []},
    //   {name: b}
    // ]
    convertObjToArr (obj) {
      const vm = this
      const res = []

      for (const i in obj) {
        const sub = {
          id: vm.index_id,
          name: i
        }

        vm.index_id += 1

        const cur = obj[i]

        if (cur.children === {}) continue

        if (obj[i].children !== {}) {
          sub.children = this.convertObjToArr(obj[i].children)
        }

        res.push(sub)
      }
      return res
    },
    // this look like sxxt， but it work, fix later
    // completeKey get complete key from number
    completeKey (number) {
      if (number === 1 || number === undefined) {
        return '/'
      }
      const vm = this
      let key = ''
      let cur = vm.generateConfigTree[0].children

      let i = 0
      while (cur.length) {
        if (cur[i].id === number) { // got number
          key += ('/' + cur[i].name)
          break
        }

        if (cur[i].id < number) {
          if (i + 1 < cur.length && cur[i + 1].id <= number) {
            i += 1 // next
          } else { // get into children
            key += ('/' + cur[i].name)
            cur = cur[i].children
            i = 0
          }
        }

        // this will not reach cur[i].id > number
      }
      return key
    }
  },
  watch: {
    choseNode (val, old) {
      if (!val.length) { // cancel select => val[]
        return
      }

      const completeKey = this.completeKey(val[0])
      console.log(completeKey)
    }
  }
}
</script>

<style lang="scss" scoped>
.tree {
  height: calc(100vh - 124px);
  overflow: auto;
}
</style>
