<template>
  <v-container
    id="auth-login"
    theme="login"
    fluid
    fill-height
    class="auth-login wallpaper background-welcome pa-6"
    :style="wallpaper()"
  >
    <v-theme-provider theme="login">
      <v-row id="auth-layout" class="auth-layout">
        <v-col cols="12" sm="9" md="6" lg="5" xl="3">
          <v-form ref="form" class="auth-login-form" accept-charset="UTF-8" @submit.prevent="onLogin">
            <v-card id="auth-login-box" class="elevation-12 auth-login-box pa-1 blur-7">
              <v-card-text>
                <p-auth-header></p-auth-header>
                <v-spacer></v-spacer>
                <v-row align="start" dense>
                  <template v-if="enterCode">
                    <v-col cols="12" class="text-body-2 text-center pb-1">
                      {{ $gettext(`Enter the code generated by your authenticator app:`) }}
                    </v-col>
                    <v-col cols="12" class="pb-0 d-flex justify-center">
                      <v-text-field
                        v-if="useRecoveryCode"
                        id="one-time-code"
                        ref="code"
                        v-model="code"
                        :disabled="loading"
                        name="code"
                        variant="solo"
                        density="comfortable"
                        type="text"
                        :placeholder="$gettext('Recovery Code')"
                        inputmode="text"
                        hide-details
                        autofocus
                        autocorrect="off"
                        autocapitalize="none"
                        autocomplete="none"
                        class="input-code text-selectable ma-4"
                        prepend-inner-icon="mdi-shield-check"
                        @keyup.enter="onLogin"
                      ></v-text-field>
                      <v-otp-input
                        v-else
                        v-model="code"
                        :disabled="loading"
                        :length="6"
                        :max-width="320"
                        variant="solo-filled"
                        base-color="surface"
                        name="one-time-code"
                        type="number"
                        :label="$gettext('Verification Code')"
                        class="input-code"
                        @keyup.enter="onLogin"
                      >
                      </v-otp-input>
                    </v-col>
                  </template>
                  <template v-else>
                    <v-col cols="12" class="pb-1">
                      <v-text-field
                        id="auth-username"
                        v-model="username"
                        :disabled="loading || enterCode"
                        name="username"
                        variant="solo"
                        density="comfortable"
                        type="text"
                        :placeholder="$gettext('Name')"
                        hide-details
                        autocorrect="off"
                        autocapitalize="none"
                        autocomplete="username"
                        class="input-username text-selectable"
                        prepend-inner-icon="mdi-account"
                        @keyup.enter="onLogin"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" class="pb-1">
                      <v-text-field
                        id="auth-password"
                        v-model="password"
                        :disabled="loading"
                        name="password"
                        variant="solo"
                        density="comfortable"
                        :type="showPassword ? 'text' : 'password'"
                        :placeholder="$gettext('Password')"
                        hide-details
                        autocorrect="off"
                        autocapitalize="none"
                        autocomplete="current-password"
                        class="input-password text-selectable"
                        :append-inner-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
                        prepend-inner-icon="mdi-lock"
                        @click:append-inner="showPassword = !showPassword"
                        @keyup.enter="onLogin"
                      ></v-text-field>
                    </v-col>
                  </template>
                  <v-col cols="12" class="auth-actions">
                    <div class="action-buttons auth-buttons pb-1 d-flex ga-3 align-center justify-center">
                      <v-btn
                        v-if="enterCode"
                        color="highlight"
                        variant="outlined"
                        :block="$vuetify.display.xs"
                        class="action-cancel opacity-80"
                        @click.stop.prevent="onCancel"
                      >
                        {{ $gettext(`Cancel`) }}
                      </v-btn>
                      <v-btn
                        v-else-if="registerUri"
                        color="highlight"
                        variant="outlined"
                        :block="$vuetify.display.xs"
                        class="action-register opacity-80"
                        @click.stop.prevent="onRegister"
                      >
                        {{ $gettext(`Create Account`) }}
                      </v-btn>
                      <v-btn
                        color="highlight"
                        variant="flat"
                        :disabled="loginDisabled"
                        :block="$vuetify.display.xs"
                        class="action-confirm"
                        @click.stop.prevent="onLogin"
                      >
                        {{ $gettext(`Sign in`) }}
                        <v-icon :icon="$config.isRtl() ? 'mdi-chevron-left' : 'mdi-chevron-right'" end></v-icon>
                      </v-btn>
                    </div>
                    <div
                      v-if="enterCode"
                      :class="{ clickable: !useRecoveryCode }"
                      class="auth-links text-center opacity-80"
                      @click.stop.prevent="onUseRecoveryCode"
                    >
                      {{ $gettext(`Can't access your authenticator app or device?`) }}
                      {{ $gettext(`Use your recovery code or contact an administrator for help.`) }}
                    </div>
                    <div v-else-if="passwordResetUri" class="auth-links text-center opacity-80">
                      <a :href="passwordResetUri" class="text-secondary">
                        {{ $gettext(`Forgot password?`) }}
                      </a>
                    </div>
                  </v-col>
                  <template v-if="config.ext.oidc.enabled && !enterCode">
                    <v-col cols="12" class="oidc-actions">
                      <v-divider />
                      <div class="text-center oidc-buttons mt-6">
                        <v-btn
                          color="highlight"
                          variant="flat"
                          :disabled="loading"
                          block
                          class="action-oidc-login"
                          @click.stop.prevent="onOidcLogin"
                        >
                          <img alt="" class="oidc-icon v-icon--start mx-1" :src="config.ext.oidc.icon" />
                          {{ $gettext(`Continue with %{provider}`, { provider: config.ext.oidc.provider }) }}
                        </v-btn>
                      </div>
                    </v-col>
                  </template>
                </v-row>
              </v-card-text>
            </v-card>
          </v-form>
        </v-col>
      </v-row>
      <p-auth-footer></p-auth-footer>
    </v-theme-provider>
  </v-container>
</template>

<script>
export default {
  name: "PPageLogin",
  data() {
    return {
      loading: false,
      username: "",
      password: "",
      showPassword: false,
      useRecoveryCode: false,
      code: "",
      enterCode: false,
      sponsor: this.$config.isSponsor(),
      config: this.$config.values,
      siteDescription: this.$config.getSiteDescription(),
      nextUrl: this.$route.params.nextUrl ? this.$route.params.nextUrl : "/",
      wallpaperUri: this.$config.values.wallpaperUri,
      registerUri: this.$config.values.registerUri,
      passwordResetUri: this.$config.values.passwordResetUri,
      rtl: this.$config.isRtl(),
    };
  },
  computed: {
    loginDisabled() {
      if (this.loading) {
        return true;
      } else if (this.enterCode) {
        const len = this.code.trim().length;
        return len !== 6 && len !== 12;
      }

      return this.username.trim() === "" || this.password.trim() === "";
    },
  },
  created() {
    this.$modal.enter(this.$isMobile);
    const authError = window.localStorage.getItem("authError");
    if (authError) {
      this.$notify.error(authError);
      window.localStorage.removeItem("authError");
    }
  },
  unmounted() {
    this.$modal.leave();
  },
  methods: {
    wallpaper() {
      if (this.wallpaperUri) {
        return `background-image: url(${this.wallpaperUri});`;
      }

      return "";
    },
    load() {
      this.$notify.blockUI();

      let route = this.$router.resolve({
        name: this.$session.getHome(),
      });

      setTimeout(() => {
        window.location = route.href;
      }, 100);
    },
    reset() {
      this.username = "";
      this.password = "";
      this.showPassword = false;
      this.useRecoveryCode = false;
      this.code = "";
      this.enterCode = false;
    },
    onCancel() {
      if (this.loading) {
        return;
      }
      this.reset();
    },
    onRegister() {
      window.location = this.registerUri;
    },
    onUseRecoveryCode() {
      if (!this.useRecoveryCode) {
        this.useRecoveryCode = true;
      }
    },
    onLogin() {
      const username = this.username.trim();
      const password = this.password.trim();
      const code = this.code.trim();

      if (username === "" || password === "") {
        return;
      }

      this.loading = true;
      this.$session
        .login(username, password, code)
        .then(() => {
          this.load();
        })
        .catch((e) => {
          if (e.response?.data?.code === 32) {
            this.enterCode = true;
            this.$nextTick(() => this.$refs.code.focus());
          }
          this.loading = false;
        });
    },
    onOidcLogin() {
      if (this.loading) {
        return;
      }

      if (this.config.ext?.oidc?.loginUri) {
        this.loading = true;
        this.$nextTick(() => {
          window.location = this.config.ext.oidc.loginUri;
        });
      } else {
        this.$notify.warn(this.$gettext("Missing or invalid configuration"));
      }
    },
  },
};
</script>
