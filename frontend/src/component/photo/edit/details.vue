<template>
  <div class="p-tab p-tab-photo-details">
    <v-form
      ref="form"
      validate-on="invalid-input"
      class="p-form p-form-photo-details-meta"
      accept-charset="UTF-8"
      @submit.prevent="save"
    >
      <div class="form-body">
        <div class="form-controls">
          <v-row dense align="start">
            <v-col cols="3" sm="2" class="form-thumb">
              <div>
                <img
                  :alt="model.Title"
                  :src="model.thumbnailUrl('tile_500')"
                  class="clickable"
                  @click.stop.prevent.exact="openPhoto()"
                />
              </div>
            </v-col>
            <v-col cols="9" sm="10" class="d-flex align-self-stretch flex-column ga-4">
              <v-text-field
                v-model="model.Title"
                :append-inner-icon="model.TitleSrc === 'manual' ? 'mdi-check' : ''"
                :disabled="disabled"
                :rules="[textRule]"
                hide-details
                :label="$pgettext('Photo', 'Title')"
                placeholder=""
                autocomplete="off"
                density="comfortable"
                class="input-title"
              ></v-text-field>
              <v-textarea
                v-model="model.Caption"
                :append-inner-icon="model.CaptionSrc === 'manual' ? 'mdi-check' : ''"
                :disabled="disabled"
                hide-details
                autocomplete="off"
                auto-grow
                :label="$gettext('Caption')"
                placeholder=""
                :rows="1"
                density="comfortable"
                class="input-caption"
              ></v-textarea>
            </v-col>
          </v-row>
          <v-row dense>
            <v-col cols="4" lg="2">
              <v-autocomplete
                v-model="model.Day"
                :disabled="disabled"
                :error="invalidDate"
                :label="$gettext('Day')"
                autocomplete="off"
                hide-details
                hide-no-data
                :items="options.Days()"
                item-title="text"
                item-value="value"
                density="comfortable"
                validate-on="input"
                :rules="rules.day(true)"
                class="input-day"
                @update:model-value="syncTime"
              >
              </v-autocomplete>
            </v-col>
            <v-col cols="4" lg="2">
              <v-autocomplete
                v-model="model.Month"
                :disabled="disabled"
                :error="invalidDate"
                :label="$gettext('Month')"
                autocomplete="off"
                hide-details
                hide-no-data
                :items="options.MonthsShort()"
                item-title="text"
                item-value="value"
                density="comfortable"
                validate-on="input"
                :rules="rules.month(true)"
                class="input-month"
                @update:model-value="syncTime"
              >
              </v-autocomplete>
            </v-col>
            <v-col cols="4" lg="2">
              <v-autocomplete
                v-model="model.Year"
                :disabled="disabled"
                :error="invalidDate"
                :label="$gettext('Year')"
                autocomplete="off"
                hide-details
                hide-no-data
                :items="options.Years(1000)"
                item-title="text"
                item-value="value"
                density="comfortable"
                validate-on="input"
                :rules="rules.year(true, 1000)"
                class="input-year"
                @update:model-value="syncTime"
              >
              </v-autocomplete>
            </v-col>
            <v-col cols="6" lg="2">
              <v-text-field
                v-model="time"
                :append-inner-icon="model.TakenSrc === 'manual' ? 'mdi-check' : ''"
                :disabled="disabled"
                :label="model.timeIsUTC() ? $gettext('Time UTC') : $gettext('Local Time')"
                autocomplete="off"
                autocorrect="off"
                autocapitalize="none"
                hide-details
                density="comfortable"
                validate-on="input"
                :rules="rules.time()"
                class="input-local-time"
                @update:model-value="setTime"
              ></v-text-field>
            </v-col>
            <v-col cols="6" lg="4">
              <v-autocomplete
                v-model="model.TimeZone"
                :disabled="disabled"
                :label="$gettext('Time Zone')"
                hide-no-data
                item-value="ID"
                item-title="Name"
                :items="options.TimeZones()"
                density="comfortable"
                class="input-timezone"
                @update:model-value="syncTime"
              ></v-autocomplete>
            </v-col>
            <v-col cols="12" sm="8" md="4">
              <v-autocomplete
                v-model="model.Country"
                :append-inner-icon="model.PlaceSrc === 'manual' ? 'mdi-check' : ''"
                :disabled="disabled"
                :readonly="!!(model.Lat || model.Lng)"
                :placeholder="$gettext('Country')"
                hide-details
                hide-no-data
                autocomplete="off"
                item-value="Code"
                item-title="Name"
                :items="countries"
                prepend-inner-icon="mdi-map-marker"
                density="comfortable"
                validate-on="input"
                :rules="rules.country(true)"
                class="input-country"
              >
              </v-autocomplete>
            </v-col>
            <v-col cols="4" md="2">
              <v-text-field
                v-model="model.Altitude"
                :disabled="disabled"
                hide-details
                flat
                autocomplete="off"
                autocorrect="off"
                autocapitalize="none"
                :label="$gettext('Altitude (m)')"
                placeholder=""
                color="surface-variant"
                density="comfortable"
                validate-on="input"
                :rules="rules.number(false, -10000, 1000000)"
                class="input-altitude"
              ></v-text-field>
            </v-col>
            <v-col cols="4" sm="6" md="3">
              <v-text-field
                v-model="model.Lat"
                :append-inner-icon="model.PlaceSrc === 'manual' ? 'mdi-check' : ''"
                :disabled="disabled"
                hide-details
                autocomplete="off"
                autocorrect="off"
                autocapitalize="none"
                :label="$gettext('Latitude')"
                placeholder=""
                density="comfortable"
                validate-on="input"
                :rules="rules.lat(false)"
                class="input-latitude"
                @paste="pastePosition"
              ></v-text-field>
            </v-col>
            <v-col cols="4" sm="6" md="3">
              <v-text-field
                v-model="model.Lng"
                :append-inner-icon="model.PlaceSrc === 'manual' ? 'mdi-check' : ''"
                :disabled="disabled"
                hide-details
                autocomplete="off"
                autocorrect="off"
                autocapitalize="none"
                :label="$gettext('Longitude')"
                placeholder=""
                density="comfortable"
                validate-on="input"
                :rules="rules.lng(false)"
                class="input-longitude"
                @paste="pastePosition"
              ></v-text-field>
            </v-col>
            <v-col cols="12" md="6" class="p-camera-select">
              <v-select
                v-model="model.CameraID"
                :append-inner-icon="model.CameraSrc === 'manual' ? 'mdi-check' : ''"
                :disabled="disabled"
                :placeholder="$gettext('Camera')"
                :menu-props="{ maxHeight: 346 }"
                autocomplete="off"
                hide-details
                item-value="ID"
                item-title="Name"
                :items="cameraOptions"
                prepend-inner-icon="mdi-camera"
                density="comfortable"
                class="input-camera"
              >
              </v-select>
            </v-col>
            <v-col cols="6" md="3">
              <v-text-field
                v-model="model.Iso"
                :disabled="disabled"
                hide-details
                autocomplete="off"
                autocorrect="off"
                autocapitalize="none"
                label="ISO"
                placeholder=""
                density="comfortable"
                validate-on="input"
                :rules="rules.number(false, 0, 1048576)"
                class="input-iso"
              ></v-text-field>
            </v-col>
            <v-col cols="6" md="3">
              <v-text-field
                v-model="model.Exposure"
                :disabled="disabled"
                hide-details
                autocomplete="off"
                autocorrect="off"
                autocapitalize="none"
                :label="$gettext('Exposure')"
                placeholder=""
                density="comfortable"
                validate-on="input"
                :rules="rules.text(false, 0, 64)"
                class="input-exposure"
              ></v-text-field>
            </v-col>
            <v-col cols="12" md="6" class="p-lens-select">
              <v-select
                v-model="model.LensID"
                :append-inner-icon="model.CameraSrc === 'manual' ? 'mdi-check' : ''"
                :disabled="disabled"
                :placeholder="$gettext('Lens')"
                :menu-props="{ maxHeight: 346 }"
                autocomplete="off"
                hide-details
                item-value="ID"
                item-title="Name"
                :items="lensOptions"
                prepend-inner-icon="mdi-camera-iris"
                density="comfortable"
                class="input-lens"
              >
              </v-select>
            </v-col>
            <v-col cols="6" md="3">
              <v-text-field
                v-model="model.FNumber"
                :disabled="disabled"
                hide-details
                autocomplete="off"
                autocorrect="off"
                autocapitalize="none"
                :label="$gettext('F Number')"
                placeholder=""
                density="comfortable"
                validate-on="input"
                :rules="rules.number(false, 0, 1048576)"
                class="input-fnumber"
              ></v-text-field>
            </v-col>
            <v-col cols="6" md="3">
              <v-text-field
                v-model="model.FocalLength"
                :disabled="disabled"
                hide-details
                autocomplete="off"
                :label="$gettext('Focal Length')"
                placeholder=""
                density="comfortable"
                validate-on="input"
                :rules="rules.number(false, 0, 1048576)"
                class="input-focal-length"
              ></v-text-field>
            </v-col>
          </v-row>
          <v-row dense>
            <v-col cols="12" md="6">
              <v-textarea
                v-model="model.Details.Subject"
                :append-inner-icon="model.Details.SubjectSrc === 'manual' ? 'mdi-check' : ''"
                :disabled="disabled"
                :rules="[textRule]"
                hide-details
                autocomplete="off"
                auto-grow
                :label="$gettext('Subject')"
                placeholder=""
                :rows="1"
                density="comfortable"
                class="input-subject"
              ></v-textarea>
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="model.Details.Copyright"
                :append-inner-icon="model.Details.CopyrightSrc === 'manual' ? 'mdi-check' : ''"
                :disabled="disabled"
                :rules="[textRule]"
                hide-details
                autocomplete="off"
                :label="$gettext('Copyright')"
                placeholder=""
                density="comfortable"
                class="input-copyright"
              ></v-text-field>
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="model.Details.Artist"
                :append-inner-icon="model.Details.ArtistSrc === 'manual' ? 'mdi-check' : ''"
                :disabled="disabled"
                :rules="[textRule]"
                hide-details
                autocomplete="off"
                :label="$gettext('Artist')"
                placeholder=""
                density="comfortable"
                class="input-artist"
              ></v-text-field>
            </v-col>
            <v-col cols="12" md="6">
              <v-textarea
                v-model="model.Details.License"
                :append-inner-icon="model.Details.LicenseSrc === 'manual' ? 'mdi-check' : ''"
                :disabled="disabled"
                :rules="[textRule]"
                hide-details
                autocomplete="off"
                auto-grow
                :label="$gettext('License')"
                placeholder=""
                :rows="1"
                density="comfortable"
                class="input-license"
              ></v-textarea>
            </v-col>
            <v-col cols="12" md="8">
              <v-textarea
                v-model="model.Details.Keywords"
                :append-inner-icon="model.Details.KeywordsSrc === 'manual' ? 'mdi-check' : ''"
                :disabled="disabled"
                hide-details
                autocomplete="off"
                auto-grow
                :label="$gettext('Keywords')"
                placeholder=""
                :rows="1"
                density="default"
                class="input-keywords"
              ></v-textarea>
            </v-col>
            <v-col cols="12" md="4">
              <v-textarea
                v-model="model.Details.Notes"
                :append-inner-icon="model.Details.NotesSrc === 'manual' ? 'mdi-check' : ''"
                :disabled="disabled"
                hide-details
                autocomplete="off"
                auto-grow
                :label="$gettext('Notes')"
                placeholder=""
                :rows="1"
                density="default"
                class="input-notes"
              ></v-textarea>
            </v-col>
          </v-row>
        </div>
      </div>
      <div v-if="!disabled" class="form-actions form-actions--sticky">
        <div class="action-buttons">
          <v-btn color="button" variant="flat" class="action-close" @click.stop="close">
            {{ $gettext(`Close`) }}
          </v-btn>
          <v-btn
            color="highlight"
            variant="flat"
            :disabled="!model?.wasChanged() && !inReview"
            class="action-apply action-approve"
            @click.stop="save(false)"
          >
            <span v-if="inReview">{{ $gettext(`Approve`) }}</span>
            <span v-else>{{ $gettext(`Apply`) }}</span>
          </v-btn>
        </div>
      </div>
    </v-form>
  </div>
</template>

<script>
import countries from "options/countries.json";
import Thumb from "model/thumb";
import Photo from "model/photo";
import * as options from "options/options";
import { rules } from "common/form";

export default {
  name: "PTabPhotoDetails",
  props: {
    model: {
      type: Object,
      default: () => new Photo(false),
    },
    uid: {
      type: String,
      default: "",
    },
  },
  data() {
    return {
      disabled: !this.$config.feature("edit"),
      config: this.$config.values,
      all: {
        colors: [{ label: this.$gettext("Unknown"), name: "" }],
      },
      readonly: this.$config.get("readonly"),
      options,
      rules,
      countries,
      featReview: this.$config.feature("review"),
      showDatePicker: false,
      showTimePicker: false,
      invalidDate: false,
      time: "",
      textRule: (v) => v.length <= this.$config.get("clip") || this.$gettext("Text too long"),
      rtl: this.$rtl,
    };
  },
  computed: {
    cameraOptions() {
      return this.config.cameras;
    },
    lensOptions() {
      return this.config.lenses;
    },
    inReview() {
      return this.featReview && this.model.Quality < 3;
    },
  },
  watch: {
    model() {
      this.syncTime();
    },
    uid() {
      this.syncTime();
    },
  },
  created() {
    this.syncTime();
  },
  methods: {
    setTime() {
      if (this.rules.isTime(this.time)) {
        this.updateModel();
      }
    },
    syncTime() {
      if (!this.model.hasId()) {
        return;
      }

      const taken = this.model.getDateTime();
      this.time = taken.toFormat("HH:mm:ss");
    },
    pastePosition(event) {
      // Auto-fills the lat and lng fields if the text in the clipboard contains two float values.
      const clipboard = event.clipboardData ? event.clipboardData : window.clipboardData;

      if (!clipboard) {
        return;
      }

      // Get values from browser clipboard.
      const text = clipboard.getData("text");

      // Trim spaces before splitting by whitespace and/or commas.
      const val = text.trim().split(/[ ,]+/);

      // Two values found?
      if (val.length >= 2) {
        // Parse values.
        const lat = parseFloat(val[0]);
        const lng = parseFloat(val[1]);

        // Lat and long must be valid floating point numbers.
        if (!isNaN(lat) && lat >= -90 && lat <= 90 && !isNaN(lng) && lng >= -180 && lng <= 180) {
          // Update model values.
          this.model.Lat = lat;
          this.model.Lng = lng;
          // Prevent default action.
          event.preventDefault();
        }
      }
    },
    updateModel() {
      if (!this.model.hasId()) {
        return;
      }

      let localDate = this.model.localDate(this.time);

      this.invalidDate = !localDate.isValid;

      if (this.invalidDate) {
        return;
      }

      if (this.model.Day === 0) {
        this.model.Day = parseInt(localDate.toFormat("d"));
      }

      if (this.model.Month === 0) {
        this.model.Month = parseInt(localDate.toFormat("L"));
      }

      if (this.model.Year === 0) {
        this.model.Year = parseInt(localDate.toFormat("y"));
      }

      const isoTime =
        localDate.toISO({
          suppressMilliseconds: true,
          includeOffset: false,
        }) + "Z";

      this.model.TakenAtLocal = isoTime;

      if (this.model.currentTimeZoneUTC()) {
        this.model.TakenAt = isoTime;
      }
    },
    left() {
      this.$emit("next");
    },
    right() {
      this.$emit("prev");
    },
    openPhoto() {
      this.$root.$refs.viewer.showThumbs(Thumb.fromFiles([this.model]), 0);
    },
    save(close) {
      if (this.invalidDate) {
        this.$notify.error(this.$gettext("Invalid date"));
        return;
      }

      this.updateModel();

      this.model.update().then(() => {
        if (close) {
          this.$emit("close");
        }

        this.syncTime();
      });
    },
    close() {
      this.$emit("close");
    },
  },
};
</script>
