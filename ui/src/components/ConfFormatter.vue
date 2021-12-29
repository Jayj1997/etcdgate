<!--
 * @Author       : jayj
 * @Date         : 2021-12-15 10:20:12
 * @Description  :
 * @LastEditors  : jayj
 * @LastEditTime : 2021-12-15 15:20:10
-->
<template>
  <v-card tile>
    <v-sheet class="sheet pa-4 black lighten-2">
      <v-select
        class="white select"
        :items="forms"
        solo
        v-model="form"
        hide-details
        label="formatter"
        :change="changeForm()"
        dense
      ></v-select>
    </v-sheet>
    <v-card-text class="formatter pa-0">
      <codemirror
        ref="cm"
        v-model="code"
        :options="cmOptions"
        class="codemirror"
      />
    </v-card-text>
  </v-card>
</template>

<script>
import { codemirror } from 'vue-codemirror'
import 'codemirror/lib/codemirror.css'
import 'codemirror/theme/base16-dark.css'
// js
import 'codemirror/mode/javascript/javascript.js'

// yaml
import 'codemirror/mode/yaml/yaml.js'

// toml
import 'codemirror/mode/toml/toml.js'

// html
import 'codemirror/mode/htmlmixed/htmlmixed.js'
// xml
import 'codemirror/mode/xml/xml.js'

export default {
  name: 'ConfFormatter',
  props: {
    code: String
  },
  data () {
    return {
      cmOptions: {
        tabSize: 2,
        mode: 'text',
        theme: 'base16-dark',
        lineNumbers: true,
        line: true,
        lineWiseCopyCut: true, // cut copy line
        matchBrackets: true
      },
      editorModes: {
        text: 'text',
        json: 'text/javascript',
        yaml: 'text/x-yaml',
        toml: 'text/x-toml',
        html: 'text/html',
        xml: 'text/xml'
      },
      oldForm: 'text',
      form: 'text',
      forms: ['text', 'json', 'yaml', 'toml', 'html', 'xml']
    }
  },
  components: {
    codemirror
  },
  methods: {
    changeForm () {
      const vm = this
      if (vm.form !== vm.oldForm) {
        this.$refs.cm.options.mode = vm.editorModes[vm.form]

        vm.oldForm = vm.form
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.sheet {
  overflow: hidden;
}
.select {
  width: 150px;
  float: right;
}
.codemirror ::v-deep .CodeMirror {
  height: calc(100vh - 76px);
}
</style>
