<template>
  <div class="p-tab p-tab-photo-people">
    <div class="pa-2 p-faces">
      <v-alert
        v-if="markers.length === 0"
        color="surface-variant"
        icon="mdi-lightbulb-outline"
        class="no-results ma-2 opacity-70"
        variant="outlined"
      >
        <div class="font-weight-bold">
          {{ $gettext(`No people found`) }}
        </div>
        <div class="mt-2">
          {{ $gettext(`You may rescan your library to find additional faces.`) }}
          {{ $gettext(`Recognition starts after indexing has been completed.`) }}
        </div>
      </v-alert>
      <div v-else class="v-row search-results face-results cards-view d-flex">
        <div v-for="m in markers" :key="m.UID" class="v-col-12 v-col-sm-6 v-col-md-4 v-col-lg-3 d-flex">
          <v-card :data-id="m.UID" :class="m.classes()" class="result not-selectable flex-grow-1">
            <v-img :src="m.thumbnailUrl('tile_320')" :transition="false" aspect-ratio="1" class="card">
              <v-btn
                v-if="!m.SubjUID && !m.Invalid"
                :ripple="false"
                class="input-reject"
                icon
                variant="text"
                density="comfortable"
                position="absolute"
                :title="$gettext('Remove')"
                @click.stop.prevent="onReject(m)"
              >
                <v-icon class="action-reject">mdi-close</v-icon>
              </v-btn>
            </v-img>
            <v-card-actions class="meta pa-0">
              <v-btn
                v-if="m.Invalid"
                :disabled="busy"
                size="large"
                variant="flat"
                block
                :rounded="false"
                class="action-undo text-center"
                :title="$gettext('Undo')"
                @click.stop="onApprove(m)"
              >
                <v-icon>mdi-undo</v-icon>
              </v-btn>
              <v-text-field
                v-else-if="m.SubjUID"
                v-model="m.Name"
                :rules="[textRule]"
                :disabled="busy"
                :readonly="true"
                autocomplete="off"
                autocorrect="off"
                hide-details
                single-line
                clearable
                persistent-clear
                clear-icon="mdi-eject"
                density="comfortable"
                class="input-name pa-0 ma-0"
                @click:clear="onClearSubject(m)"
              ></v-text-field>
              <v-combobox
                v-else
                v-model:search="m.Name"
                :items="$config.values.people"
                item-title="Name"
                item-value="Name"
                :disabled="busy"
                return-object
                hide-no-data
                :menu-props="menuProps"
                hide-details
                single-line
                open-on-clear
                focused
                append-icon=""
                prepend-inner-icon="mdi-account-plus"
                density="comfortable"
                class="input-name pa-0 ma-0"
                @blur="onSetName(m)"
                @update:model-value="(person) => onSetPerson(m, person)"
                @keyup.enter.native="onSetName(m)"
              >
              </v-combobox>
            </v-card-actions>
          </v-card>
        </div>
      </div>
    </div>
    <p-confirm-action
      :show="confirm.show"
      icon="mdi-account-plus"
      :icon-size="42"
      :text="confirm?.model?.Name ? $gettext('Add %{s}?', { s: confirm.model.Name }) : $gettext('Add person?')"
      @close="onCancelSetName"
      @confirm="onConfirmSetName"
    ></p-confirm-action>
  </div>
</template>

<script>
import Marker from "model/marker";
import PConfirmAction from "component/confirm/action.vue";

export default {
  name: "PTabPhotoPeople",
  components: { PConfirmAction },
  props: {
    model: {
      type: Object,
      default: () => {},
    },
    uid: {
      type: String,
      default: "",
    },
  },
  data() {
    return {
      busy: false,
      markers: this.model.getMarkers(true),
      imageUrl: this.model.thumbnailUrl("fit_720"),
      disabled: !this.$config.feature("edit"),
      config: this.$config.values,
      readonly: this.$config.get("readonly"),
      confirm: {
        show: false,
        model: new Marker(),
        text: this.$gettext("Add person?"),
      },
      menuProps: {
        closeOnClick: false,
        closeOnContentClick: true,
        openOnClick: false,
        density: "compact",
        maxHeight: 300,
      },
      textRule: (v) => {
        if (!v || !v.length) {
          return this.$gettext("Name");
        }

        return v.length <= this.$config.get("clip") || this.$gettext("Name too long");
      },
    };
  },
  watch: {
    model: function () {
      this.refresh();
    },
  },
  methods: {
    refresh() {
      this.markers = this.model.getMarkers(true);
      this.imageUrl = this.model.thumbnailUrl("fit_720");
    },
    onReject(model) {
      if (this.busy || !model) return;

      this.busy = true;
      this.$notify.blockUI();

      model.reject().finally(() => {
        this.$notify.unblockUI();
        this.busy = false;
      });
    },
    onApprove(model) {
      if (this.busy || !model) return;

      this.busy = true;

      model.approve().finally(() => (this.busy = false));
    },
    onClearSubject(model) {
      if (this.busy || !model) return;

      this.busy = true;
      this.$notify.blockUI();

      model.clearSubject(model).finally(() => {
        this.$notify.unblockUI();
        this.busy = false;
      });
    },
    onSetPerson(model, person) {
      if (typeof person === "object" && model?.UID && person?.UID && person?.Name) {
        model.Name = person.Name;
        model.SubjUID = person.UID;
        this.setName(model);
      }

      return true;
    },
    onSetName(model) {
      if (this.busy || !model) {
        return;
      }

      const name = model?.Name;

      if (!name) {
        this.onCancelSetName();
        return;
      }

      this.confirm.model = model;

      const people = this.$config.values?.people;

      if (people) {
        const found = people.find((person) => person.Name.localeCompare(name, "en", { sensitivity: "base" }) === 0);
        if (found) {
          model.Name = found.Name;
          model.SubjUID = found.UID;
          this.setName(model);
          return;
        }
      }

      model.Name = name;
      model.SubjUID = "";
      this.confirm.show = true;
    },
    onConfirmSetName() {
      if (!this.confirm?.model?.Name) {
        return;
      }

      this.setName(this.confirm.model);
    },
    onCancelSetName() {
      this.confirm.show = false;
    },
    setName(model) {
      if (this.busy || !model) {
        return;
      }

      this.busy = true;
      this.$notify.blockUI();

      return model.setName().finally(() => {
        this.$notify.unblockUI();
        this.busy = false;
        this.confirm.model = null;
        this.confirm.show = false;
      });
    },
  },
};
</script>
