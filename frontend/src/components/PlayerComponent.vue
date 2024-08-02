<script lang="ts">
import {CurrentTrack, NextMusic, PauseMusic, PlayMusic} from "app/wailsjs/go/main/App";

export default {
  data: () => ({
    currentTrackName: "",
  }),
  methods: {
    play() {
      PlayMusic().then(() => this.getCurrentTrack())
    },
    next() {
      NextMusic().then(() => this.getCurrentTrack())
    },
    pause() {
      PauseMusic().then(() => this.getCurrentTrack())
    },
    getCurrentTrack() {
      CurrentTrack().then((r) => {
        this.currentTrackName = r.Label
      })
    }
  }
}
</script>

<template>
  <div class="q-pa-md example-row-equal-width">
    <div class="row">
      <div class="col">
        <q-input
          readonly
          color="black"
          v-model="currentTrackName"
        />
      </div>
    </div>
    <div class="row justify-center">
      <div class="col-3 self-start">
        <q-btn round icon="play_arrow" @click="play"></q-btn>
      </div>
      <div class="col-3 self-center">
        <q-btn round icon="skip_next" @click="next"></q-btn>
      </div>
      <div class="col-3 self-end">
        <q-btn size="md" round icon="pause" @click="pause"></q-btn>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>
