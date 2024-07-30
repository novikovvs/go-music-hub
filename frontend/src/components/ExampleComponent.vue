<template>
  <div>
    <q-btn v-if="!addr && !loading" color="white" text-color="black" label="Start" @click="setup"/>
    <q-skeleton v-if="loading" type="QBtn"/>

    <q-card v-if="addr && !loading" flat bordered class="my-card">
      <q-card-section>
        <div class="row items-center no-wrap">
          <div class="col">
            <div class="text-h6">Адрес комнаты</div>
          </div>
        </div>
      </q-card-section>

      <q-card-section>
        <q-input outlined v-model="addr" readonly @click="copyAddr"/>
      </q-card-section>

      <q-separator/>

      <q-card-actions class=" flex-center">
        <q-btn size="md" class="flex-sm" color="red" text-color="black" label="Остановить" @click="stop"/>
      </q-card-actions>
    </q-card>
  </div>
</template>

<script setup lang="ts">
import {CreateRoom, StopSrv} from "app/wailsjs/go/main/App";
import {ref} from "vue";
</script>

<script lang="ts">
import {CreateRoom, RoomAlive, StopSrv} from "app/wailsjs/go/main/App";
import {ref} from "vue";
import {Notify} from 'quasar'


export default {
  data: () => ({
    addr: null,
    loading: false,
  }),
  mounted() {
    RoomAlive().then((r: boolean | string) => {
      this.addr = r
    })
  },
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

      StopSrv().then((r) => {
        this.loading = false

        r ? this.addr = null : this.addr = this.addr
      })
    },
    async copyAddr() {
      try {
        await navigator.clipboard.writeText(this.addr);
        Notify.create({message: 'Скопировано!', color: "green", progress: true, timeout: 1000})
      } catch ($e) {
        Notify.create({message: 'Неудалось скопировать!', color: "red", progress: true, timeout: 1000})
      }
    }
  }
}
</script>
