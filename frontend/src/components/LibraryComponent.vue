<template>
  <div class="q-ma-lg">
    <div class="q-ma-md">
      <q-btn color="primary" icon="sync" label="Обновить" size="sm" @click="sync" :loading="syncing">
        <template v-slot:loading>
          <q-spinner-gears/>
        </template>
      </q-btn>
    </div>
    <q-separator color="white"/>
    <q-tree
      v-if="libraries.length >0"
      :nodes="libraries"
      default-expand-all
      node-key="Label"
      label-key="Label"
      text-color="white"
      children-key="Children"
    >
      <template v-slot:default-body="prop">
        <div v-if="prop.node.isTrack">
          <span class="text-weight-bold" @click="$router.push('/player')"><q-icon size="md" name="play_arrow"/></span>
        </div>
      </template>

    </q-tree>
  </div>
</template>

<script lang="ts">
import {GetLibrary} from "app/wailsjs/go/main/App";
import {TrackLibrary} from "app/wailsjs/go/models";

function libraryAdapter(r: TrackLibrary[]) {
  return [{
    "Label": "Библиотеки музыки",
    "Children": r.map((value: TrackLibrary) => {
      return {
        Label: value.Label,
        Children: value.Children.map((chilValue) => {
          return {
            Label: chilValue.Label,
            Children: chilValue.Children,
            isTrack: true,
          }
        })
      }

    })
  }]
}

export default {
  data: () => ({
    libraries: [],
    syncing: false,
  }),
  mounted() {
    GetLibrary().then((r) => {
      this.libraries = libraryAdapter(r)
    })
  },
  methods: {
    sync() {
      this.syncing = true
      GetLibrary().then((r) => {
        this.libraries = libraryAdapter(r)
        this.syncing = false

      })
    }
  }
}
</script>
