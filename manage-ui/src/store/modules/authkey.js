// src/store/modules/authkey.js
import axios from 'axios';

const state = {
  authKeys: [], // AuthKey列表
};

const mutations = {
  setAuthKeys(state, authKeys) {
    state.authKeys = authKeys;
  },
};

const actions = {
  async fetchAuthKeysPaged({ commit, dispatch, rootState }, { limit }) {
    const headers = await dispatch('getHeaders', null, { root: true });
    const { data } = await axios.post(
      `${rootState.baseUrl}/api?api_id=e8181ee2`,
      { limit },
      { headers }
    );
    commit('setAuthKeys', data.data);
    return data;
  },
  async addAuthKey({ commit, state, rootState, dispatch }, authKeyData) {
    const headers = await dispatch('getHeaders', null, { root: true });
    const { data } = await axios.post(`${rootState.baseUrl}/api?api_id=81e616e3`, authKeyData, {
      headers,
    });
    if (data.code === 0) {
      commit('setAuthKeys', [...state.authKeys, data.data]);
    }
    return data;
  },
  async removeAuthKey({ commit, state, rootState, dispatch }, authKeyId) {
    const headers = await dispatch('getHeaders', null, { root: true });
    const { data } = await axios.post(`${rootState.baseUrl}/api?api_id=81acavb2`, { auth_key: authKeyId }, {
      headers,
    });
    if (data.code === 0) {
      commit('setAuthKeys', state.authKeys.filter(p => p.id !== authKeyId));
    }
    return data;
  },
  async updateAuthKey({ commit, state, rootState, dispatch }, authKeyData) {
    const headers = await dispatch('getHeaders', null, { root: true });
    const { data } = await axios.post(`${rootState.baseUrl}/api?api_id=716e71c6`, authKeyData, {
      headers,
    });
    if (data.code === 0) {
      commit('setAuthKeys', [...state.authKeys, data.data]);
    }
    return data;
  },
};

export default {
  namespaced: true,
  state,
  mutations,
  actions,
};
