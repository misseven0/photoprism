<template>
  <div class="p-page p-page-help">
    <v-toolbar flat color="secondary" :density="$vuetify.display.smAndDown ? 'compact' : 'default'">
      <v-toolbar-title class="flex-grow-1">
        {{ $gettext(`Help`) }}
      </v-toolbar-title>

      <v-btn
        icon
        href="https://docs.photoprism.app/getting-started/troubleshooting/"
        target="_blank"
        class="action-info mx-2"
        :title="$gettext('Learn more')"
      >
        <v-icon icon="mdi-book-open-page-variant" size="26" color="surface-variant"></v-icon>
      </v-btn>
    </v-toolbar>

    <div class="px-6 py-4">
      <p-help-websockets></p-help-websockets>
    </div>

    <p-about-footer></p-about-footer>
  </div>
</template>

<script>
import PHelpWebsockets from "./help/websockets.vue";
import PAboutFooter from "component/about/footer.vue";

export default {
  name: "PPageHelp",
  components: {
    PHelpWebsockets,
    PAboutFooter,
  },
  data() {
    return {
      topic: this.$route.params.pathMatch,
    };
  },
};
</script>
