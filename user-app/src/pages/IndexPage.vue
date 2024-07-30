<template>
  <q-page class="flex flex-center">
    <q-form
      @submit="onSubmit"
      class="q-gutter-md"
    >
      <q-input
        filled
        v-model="url"
        label="Ссылка на Youtube"
        @update:model-value="getMetaData"
      />


      <div v-if="thumbnailUrl">
        <q-card class="my-card q-mb-sm">
          <q-img :src="thumbnailUrl">
          </q-img>
        </q-card>
        <q-btn label="Добавить" type="submit" color="primary"/>
      </div>
    </q-form>
  </q-page>
</template>

<script>
import {api} from "src/axios/axios";

export default {
  data: () => ({
    url: null,
    thumbnailUrl: null,
    videoId: null
  }),
  methods: {
    getMetaData() {
      this.videoId = getVideoId(this.url)
      this.thumbnailUrl = getYoutubeThumbnail(this.videoId, 'max')
    },
    onSubmit() {
      api.post('/submit', {
        "video_id": this.videoId,
        "url": this.url,
      }).then((r) =>
        Notify.create({message: 'Удачно!', color: "green", progress: true, timeout: 1000})
      )
    },
  }
}

function getYoutubeThumbnail(videoId, quality) {
  if (videoId) {
    if (typeof quality == "undefined") {
      quality = 'high';
    }

    var quality_key = 'maxresdefault'; // Max quality
    if (quality === 'low') {
      quality_key = 'sddefault';
    } else if (quality === 'medium') {
      quality_key = 'mqdefault';
    } else if (quality === 'high') {
      quality_key = 'hqdefault';
    }

    return "http://img.youtube.com/vi/" + videoId + "/" + quality_key + ".jpg";
  }
  return false;
}

function getVideoId(url) {
  var result, videoId
  if (result = url.match(/(?<=v=)\w+/)) {
    videoId = result.pop();
  } else if (result = url.match(/youtu.be\/(.{11})/)) {
    videoId = result.pop();
  }

  return videoId
}
</script>
