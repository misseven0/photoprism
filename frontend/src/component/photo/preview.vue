<template>
  <div class="p-photo-preview ra-6 pa-0 ma-0 elevation-0 v-card v-sheet v-sheet--tile no-transition" :title="title">
    <div class="v-responsive v-image card elevation-0 clickable" @click.prevent.stop="openPhoto">
      <div class="v-responsive__sizer" style="padding-bottom: 100%"></div>
      <div class="v-image__image v-image__image--cover w-100" :style="cover"></div>
    </div>
  </div>
</template>
<script>
import Photo from "model/photo";
import Thumb from "model/thumb";

export default {
  name: "PPhotoPreview",
  props: {
    model: {
      type: Object,
      default: () => new Photo(false),
    },
  },
  data() {
    return {
      url: this.model.thumbnailUrl("tile_500"),
      title: this.model.Title ? this.model.Title : "",
    };
  },
  computed: {
    cover() {
      return `background-image: url('${this.url}'); background-position: center center;background-size: cover;`;
    },
  },
  watch: {
    model() {
      this.url = this.model.thumbnailUrl("tile_500");
      this.title = this.model.Title ? this.model.Title : "";
    },
  },
  methods: {
    openPhoto() {
      if (!this.$viewer || !this.model) {
        return;
      }

      this.$root.$refs.viewer.showThumbs(Thumb.fromFiles([this.model]), 0);
    },
  },
};
</script>
