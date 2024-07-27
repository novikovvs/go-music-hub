<template>
  <div>
    <q-btn v-if="!addr && !loading" color="white" text-color="black" label="Start" @click="setup"/>
    <q-btn v-if="addr && !loading" color="white" text-color="black" label="Stop" @click="stop"/>
    <q-skeleton v-if="loading" type="QBtn" />
    <p class="text-center" v-if="!loading">
      {{ addr }}
    </p>
  </div>
</template>

<script setup lang="ts">
import {CreateRoom, StopSrv} from "app/wailsjs/go/main/App";
import {ref} from "vue";

</script>

<script lang="ts">
import {CreateRoom, StopSrv} from "app/wailsjs/go/main/App";
import {ref} from "vue";

export default {
  data: () => ({
    addr: null,
    loading: false,
  }),
  methods: {
    setup() {
      this.loading = true
      CreateRoom().then((a) => {
        this.loading = false
        this.addr = a
        }
      )
    },
    stop() {
      this.loading = true

      StopSrv().then((r) =>{
        this.loading = false

        r ? this.addr = null : this.addr = this.addr
      })
    }
  }
}
</script>
