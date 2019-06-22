import Vue from 'vue'
import Vuex from 'vuex'

import interfaceModule from './modules/interface.js'
import instanceModule from './modules/instance.js'
import oauthModule from './modules/oauth.js'
import authFlowModule from './modules/auth_flow.js'

import createPersistedState from './lib/persisted_state.js'

import afterStoreSetup from './boot/after_store.js'

Vue.config.productionTip = false

Vue.use(Vuex)

const persistedStateOptions = {
  paths: [
    'config',
    'users.lastLoginName',
    'oauth'
  ]
};

(async () => {
  const persistedState = await createPersistedState(persistedStateOptions)
  const store = new Vuex.Store({
    modules: {
      interface: interfaceModule,
      instance: instanceModule,
      oauth: oauthModule,
      authFlow: authFlowModule
    },
    plugins: [persistedState]
  })

  afterStoreSetup({ store })
})()
