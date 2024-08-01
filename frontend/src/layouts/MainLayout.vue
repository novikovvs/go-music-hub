<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated class="bg-dark">
      <q-toolbar>
        <q-btn flat @click="toggleLeftDrawer" dense icon="audio_file"/>
        <q-toolbar-title>
          Это лечится музыкой
        </q-toolbar-title>
        <q-btn flat @click="$router.push('/')" dense label="Бот"/>
        <q-btn flat @click="$router.push('/player')" dense label="Проигрыватель"/>
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      :breakpoint="700"
      elevated
      class="bg-dark text-white"
    >
      <q-scroll-area class="fit">
        <library-component></library-component>
      </q-scroll-area>
    </q-drawer>


    <q-page-container>
      <router-view/>
    </q-page-container>
  </q-layout>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import LibraryComponent from "components/LibraryComponent.vue";

export default defineComponent({
  components: {LibraryComponent},
  data: () => ({
    leftDrawerOpen: false
  }),
  methods: {
    toggleLeftDrawer() {
      this.leftDrawerOpen = !this.leftDrawerOpen
    }
  }
})
</script>
