<template>
  <v-form
    ref="form"
    validate-on="invalid-input"
    autocomplete="off"
    class="p-photo-toolbar"
    accept-charset="UTF-8"
    :class="{ embedded: embedded }"
    @submit.prevent="updateQuery()"
  >
    <v-toolbar
      :density="$vuetify.display.smAndDown && !embedded ? 'compact' : 'default'"
      :height="embedded ? 45 : undefined"
      class="page-toolbar"
      color="secondary"
    >
      <template v-if="!embedded">
        <v-text-field
          :model-value="filter.q"
          hide-details
          clearable
          single-line
          overflow
          rounded
          variant="solo-filled"
          :density="density"
          validate-on="invalid-input"
          autocorrect="off"
          autocapitalize="none"
          autocomplete="off"
          :placeholder="$gettext('Search')"
          prepend-inner-icon="mdi-magnify"
          color="surface-variant"
          class="input-search background-inherit elevation-0"
          @update:modelValue="
            (v) => {
              updateFilter({ q: v });
            }
          "
          @keyup.enter="() => onUpdate()"
          @click:clear="
            () => {
              onUpdate({ q: '' });
            }
          "
        ></v-text-field>

        <v-btn
          v-if="filter.latlng"
          icon
          :title="$gettext('Show more')"
          class="action-clear-location"
          @click.stop="clearLocation()"
        >
          <v-icon>mdi-map-marker-off</v-icon>
        </v-btn>

        <v-btn icon class="hidden-xs action-reload" :title="$gettext('Reload')" @click.stop="refresh()">
          <v-icon>mdi-refresh</v-icon>
        </v-btn>

        <v-btn
          v-if="settings.view === 'list'"
          icon
          class="action-view-mosaic"
          :title="$gettext('Toggle View')"
          @click.stop="setView('mosaic')"
        >
          <v-icon>mdi-view-comfy</v-icon>
        </v-btn>
        <v-btn
          v-else-if="settings.view === 'cards' && listView"
          icon
          class="action-view-list"
          :title="$gettext('Toggle View')"
          @click.stop="setView('list')"
        >
          <v-icon>mdi-view-list</v-icon>
        </v-btn>
        <v-btn
          v-else-if="settings.view === 'cards'"
          icon
          class="action-view-mosaic"
          :title="$gettext('Toggle View')"
          @click.stop="setView('mosaic')"
        >
          <v-icon>mdi-view-comfy</v-icon>
        </v-btn>
        <v-btn v-else icon class="action-view-cards" :title="$gettext('Toggle View')" @click.stop="setView('cards')">
          <v-icon>mdi-view-column</v-icon>
        </v-btn>
        <v-btn
          v-if="canDelete && context === 'archive' && config.count.archived > 0"
          icon
          class="action-delete-all"
          :title="$gettext('Delete All')"
          @click.stop="deleteAll()"
        >
          <v-icon>mdi-delete-sweep</v-icon>
        </v-btn>
        <v-btn
          v-else-if="canUpload"
          icon
          class="hidden-sm-and-down action-upload"
          :title="$gettext('Upload')"
          @click.stop="showUpload()"
        >
          <v-icon>mdi-cloud-upload</v-icon>
        </v-btn>
        <v-btn
          :icon="expanded ? 'mdi-chevron-up' : 'mdi-chevron-down'"
          :title="$gettext('Expand Search')"
          class="action-expand"
          :class="{ 'action-expand--active': expanded }"
          @click.stop="toggleExpansionPanel"
        ></v-btn>
      </template>
      <template v-else>
        <v-spacer></v-spacer>
        <v-btn v-if="canAccessLibrary" icon :title="$gettext('Browse')" class="action-browse" @click.stop="onBrowse">
          <v-icon size="20">mdi-tab</v-icon>
        </v-btn>
        <v-btn v-if="onClose !== undefined" icon :title="$gettext('Close')" class="action-close" @click.stop="onClose">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </template>
    </v-toolbar>

    <div class="toolbar-expansion-panel">
      <v-expand-transition>
        <v-card v-show="expanded" flat color="secondary">
          <v-card-text class="dense">
            <v-row align="center" dense>
              <v-col cols="12" sm="6" md="3" class="p-countries-select">
                <v-select
                  :model-value="filter.country"
                  :label="$gettext('Country')"
                  :menu-props="{ maxHeight: 346 }"
                  single-line
                  hide-details
                  variant="solo-filled"
                  :density="density"
                  :items="countryOptions"
                  item-title="Name"
                  item-value="ID"
                  class="input-countries"
                  @update:model-value="
                    (v) => {
                      onUpdate({ country: v });
                    }
                  "
                >
                </v-select>
              </v-col>
              <v-col cols="12" sm="6" md="3" class="p-camera-select">
                <v-select
                  :model-value="filter.camera"
                  :label="$gettext('Camera')"
                  :menu-props="{ maxHeight: 346 }"
                  single-line
                  hide-details
                  variant="solo-filled"
                  :density="density"
                  :items="cameraOptions"
                  item-title="Name"
                  item-value="ID"
                  @update:model-value="
                    (v) => {
                      onUpdate({ camera: v });
                    }
                  "
                >
                </v-select>
              </v-col>
              <v-col cols="12" sm="6" md="3" class="p-view-select">
                <v-select
                  id="viewSelect"
                  :model-value="settings.view"
                  :label="$gettext('View')"
                  single-line
                  hide-details
                  variant="solo-filled"
                  :density="density"
                  :items="viewOptions"
                  item-title="text"
                  item-value="value"
                  @update:model-value="
                    (v) => {
                      setView(v);
                    }
                  "
                >
                </v-select>
              </v-col>
              <v-col cols="12" sm="6" md="3" class="p-time-select">
                <v-select
                  :model-value="filter.order"
                  :label="$gettext('Sort Order')"
                  :menu-props="{ maxHeight: 400 }"
                  single-line
                  variant="solo-filled"
                  :density="density"
                  :items="sortOptions"
                  item-title="text"
                  item-value="value"
                  @update:model-value="
                    (v) => {
                      onUpdate({ order: v });
                    }
                  "
                >
                </v-select>
              </v-col>
              <v-col cols="12" sm="6" md="3" class="p-year-select">
                <v-select
                  :model-value="filter.year"
                  :label="$gettext('Year')"
                  :menu-props="{ maxHeight: 346 }"
                  single-line
                  variant="solo-filled"
                  :density="density"
                  :items="yearOptions()"
                  item-title="text"
                  item-value="value"
                  @update:model-value="
                    (v) => {
                      onUpdate({ year: v });
                    }
                  "
                >
                </v-select>
              </v-col>
              <v-col cols="12" sm="6" md="3" class="p-month-select">
                <v-select
                  :model-value="filter.month"
                  :label="$gettext('Month')"
                  :menu-props="{ maxHeight: 346 }"
                  single-line
                  variant="solo-filled"
                  :density="density"
                  :items="monthOptions()"
                  item-title="text"
                  item-value="value"
                  @update:model-value="
                    (v) => {
                      onUpdate({ month: v });
                    }
                  "
                >
                </v-select>
              </v-col>
              <!-- v-col cols="12" sm="6" md="3" class="p-lens-select">
                <v-select @change="dropdownChange"
                          :label="labels.lens"
                          flat
                          variant="solo-filled"
                          hide-details
                          color="surface-variant"
                          bg-color="secondary-light"
                          item-value="ID"
                          item-title="Model"
                          v-model="filter.lens"
                          :items="lensOptions">
                </v-select>
            </v-col -->
              <v-col cols="12" sm="6" md="3" class="p-color-select">
                <v-select
                  :model-value="filter.color"
                  :label="$gettext('Color')"
                  :menu-props="{ maxHeight: 346 }"
                  single-line
                  hide-details
                  variant="solo-filled"
                  :density="density"
                  :items="colorOptions()"
                  item-title="Name"
                  item-value="Slug"
                  @update:model-value="
                    (v) => {
                      onUpdate({ color: v });
                    }
                  "
                >
                </v-select>
              </v-col>
              <v-col cols="12" sm="6" md="3" class="p-category-select">
                <v-select
                  :model-value="filter.label"
                  :label="$gettext('Category')"
                  :menu-props="{ maxHeight: 346 }"
                  single-line
                  hide-details
                  variant="solo-filled"
                  :density="density"
                  :items="categoryOptions"
                  item-title="Name"
                  item-value="Slug"
                  @update:model-value="
                    (v) => {
                      onUpdate({ label: v });
                    }
                  "
                >
                </v-select>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-expand-transition>
    </div>
    <p-photo-delete-dialog
      :show="dialog.delete"
      :text="$gettext('Are you sure you want to delete all archived pictures?')"
      :action="$gettext('Delete All')"
      @close="dialog.delete = false"
      @confirm="batchDelete"
    >
    </p-photo-delete-dialog>
  </v-form>
</template>
<script>
import Event from "pubsub-js";
import * as options from "options/options";
import $api from "common/api";
import $notify from "common/notify";

export default {
  name: "PPhotoToolbar",
  props: {
    context: {
      type: String,
      default: "photos",
    },
    filter: {
      type: Object,
      default: () => {},
    },
    staticFilter: {
      type: Object,
      default: () => {},
    },
    updateFilter: {
      type: Function,
      default: () => {},
    },
    updateQuery: {
      type: Function,
      default: () => {},
    },
    settings: {
      type: Object,
      default: () => {},
    },
    refresh: {
      type: Function,
      default: () => {},
    },
    onClose: {
      type: Function,
      default: undefined,
    },
    embedded: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    const features = this.$config.getSettings().features;
    const readonly = this.$config.get("readonly");

    return {
      expanded: false,
      experimental: this.$config.get("experimental"),
      isFullScreen: !!document.fullscreenElement,
      config: this.$config.values,
      readonly: readonly,
      canUpload: !readonly && !this.embedded && this.$config.allow("files", "upload") && features.upload,
      canDelete: !readonly && !this.embedded && this.$config.allow("photos", "delete") && features.delete,
      canAccessLibrary: this.$config.allow("photos", "access_library"),
      listView: this.$config.getSettings()?.search?.listView,
      all: {
        countries: [{ ID: "", Name: this.$gettext("All Countries") }],
        cameras: [{ ID: 0, Name: this.$gettext("All Cameras") }],
        lenses: [{ ID: 0, Name: this.$gettext("All Lenses") }],
        colors: [{ Slug: "", Name: this.$gettext("All Colors") }],
        categories: [{ Slug: "", Name: this.$gettext("All Categories") }],
        months: [{ value: 0, text: this.$gettext("All Months") }],
        years: [{ value: 0, text: this.$gettext("All Years") }],
      },
      dialog: {
        delete: false,
      },
    };
  },
  computed: {
    density() {
      return this.$vuetify.display.smAndDown ? "compact" : "comfortable";
    },
    countryOptions() {
      return this.all.countries.concat(this.config.countries);
    },
    cameraOptions() {
      return this.all.cameras.concat(this.config.cameras);
    },
    categoryOptions() {
      return this.all.categories.concat(this.config.categories);
    },
    viewOptions() {
      if (this.$config.getSettings()?.search?.listView) {
        return [
          { value: "mosaic", text: this.$gettext("Mosaic") },
          { value: "cards", text: this.$gettext("Cards") },
          { value: "list", text: this.$gettext("List") },
        ];
      } else {
        return [
          { value: "mosaic", text: this.$gettext("Mosaic") },
          { value: "cards", text: this.$gettext("Cards") },
        ];
      }
    },
    sortOptions() {
      switch (this.context) {
        case "archive":
          return [
            { value: "newest", text: this.$gettext("Newest First") },
            { value: "oldest", text: this.$gettext("Oldest First") },
            { value: "added", text: this.$gettext("Recently Added") },
            { value: "archived", text: this.$gettext("Recently Archived") },
            { value: "title", text: this.$gettext("Picture Title") },
            { value: "name", text: this.$gettext("File Name") },
            { value: "size", text: this.$gettext("File Size") },
            { value: "duration", text: this.$gettext("Video Duration") },
          ];
        case "hidden":
        case "review":
          return [
            { value: "newest", text: this.$gettext("Newest First") },
            { value: "oldest", text: this.$gettext("Oldest First") },
            { value: "added", text: this.$gettext("Recently Added") },
            { value: "title", text: this.$gettext("Picture Title") },
            { value: "name", text: this.$gettext("File Name") },
            { value: "size", text: this.$gettext("File Size") },
            { value: "duration", text: this.$gettext("Video Duration") },
          ];
        default:
          return [
            { value: "newest", text: this.$gettext("Newest First") },
            { value: "oldest", text: this.$gettext("Oldest First") },
            { value: "added", text: this.$gettext("Recently Added") },
            { value: "edited", text: this.$gettext("Recently Edited") },
            { value: "title", text: this.$gettext("Picture Title") },
            { value: "name", text: this.$gettext("File Name") },
            { value: "size", text: this.$gettext("File Size") },
            { value: "duration", text: this.$gettext("Video Duration") },
            { value: "similar", text: this.$gettext("Visual Similarity") },
            { value: "relevance", text: this.$gettext("Most Relevant") },
          ];
      }
    },
  },
  methods: {
    hideExpansionPanel() {
      if (this.expanded) {
        this.expanded = false;
      }
    },
    toggleExpansionPanel() {
      this.expanded = !this.expanded;
    },
    colorOptions() {
      return this.all.colors.concat(options.Colors());
    },
    monthOptions() {
      return this.all.months.concat(options.Months());
    },
    yearOptions() {
      return this.all.years.concat(options.IndexedYears());
    },
    setView(name) {
      if (name) {
        if (name === "list" && !this.listView) {
          name = "mosaic";
        }
        this.hideExpansionPanel();
        this.refresh({ view: name });
      }
    },
    showUpload() {
      Event.publish("dialog.upload");
    },
    deleteAll() {
      if (!this.canDelete) {
        return;
      }

      this.dialog.delete = true;
    },
    clearLocation() {
      this.$router.push({ name: "browse" });
    },
    onBrowse() {
      const route = { name: "places_browse", query: this.staticFilter };

      if (this.$isMobile) {
        this.$router.push(route);
      } else {
        // Open in a new tab on desktop browsers.
        const routeUrl = this.$router.resolve(route).href;

        if (routeUrl) {
          window.open(routeUrl, "_blank");
        }
      }
    },
    onUpdate(v) {
      this.updateQuery(v);
    },
    batchDelete() {
      if (!this.canDelete) {
        return;
      }

      this.dialog.delete = false;

      $api.post("batch/photos/delete", { all: true }).then(() => this.onDeleted());
    },
    onDeleted() {
      $notify.success(this.$gettext("Permanently deleted"));
      this.$clipboard.clear();
    },
  },
};
</script>
