<template>
  <div>
    <q-card flat bordered class="my-card">
      <q-card-section>
        <div class="row items-center no-wrap">
          <div class="col col-2">
            <q-btn round icon="sync" size="sm" @click="checkAlive" :loading="syncing">
              <template v-slot:loading>
                <q-spinner-gears/>
              </template>
            </q-btn>
          </div>
          <div class="col col-4">
            <div class="text-h6 block">
              Бот
              <q-badge rounded :color="botAlive ? `green` : `red`"/>
            </div>
          </div>
        </div>
      </q-card-section>

      <q-card-section>
        <q-input outlined v-model="addr" label="Апи ключ" :readonly="botAlive"/>
      </q-card-section>

      <q-separator/>

      <q-card-actions class=" flex-center">
        <q-btn v-if="!botAlive" color="white" text-color="black" label="Запустить" @click="setup" :loading="loading">
          <template v-slot:loading>
            <q-spinner-gears/>
          </template>
        </q-btn>
        <q-btn v-else size="md" class="flex-sm" color="red" text-color="black" label="Остановить" @click="stop"
               :loading="loading">
          <template v-slot:loading>
            <q-spinner-gears/>
          </template>
        </q-btn>
      </q-card-actions>
    </q-card>
  </div>
</template>

<script lang="ts">
import {CreateRoomBot, RoomAlive, StopSrv} from 'app/wailsjs/go/main/App';

export default {
  data: () => ({
    addr: "",
    loading: false,
    botAlive: false,
    syncing: false,
  }),
  mounted() {
    RoomAlive().then((r: string) => {
      if (r.length > 0) {
        this.botAlive = true
        this.addr = r
        return
      }
      this.botAlive = false
    })
  },
  methods: {
    stop() {
      this.loading = true
      StopSrv().then(() => {
        this.loading = false
        this.botAlive = false
      })
    },

    setup() {
      this.loading = true
      CreateRoomBot(this.addr).then((r) => {
        this.checkAlive()
        this.loading = false
      })
    },

    checkAlive() {
      this.syncing = true
      RoomAlive().then((r: string) => {
        this.syncing = false
        if (r.length > 0) {
          this.botAlive = true
          return
        }
        this.botAlive = false
      })
    }
  }
}
</script>
