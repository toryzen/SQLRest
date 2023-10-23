// src/store/modules/db.js
import axios from 'axios';

const state = {
  dbs: [], // 数据库列表
};

const mutations = {
  setDBs(state, dbs) {
    state.dbs = dbs;
  },
};

const actions = {
  async fetchDBsPaged({ commit, dispatch, rootState }, { limit }) {
    const headers = await dispatch('getHeaders', null, { root: true });
    const { data } = await axios.post(
      `${rootState.baseUrl}/api?api_id=e8e8181r`,
      { limit },
      { headers, }
    );
    commit('setDBs', data.data);
    return data;
  },
  async addDb({ commit, state, rootState, dispatch }, dbData) {
    const headers = await dispatch('getHeaders', null, { root: true });
    const { data } = await axios.post(`${rootState.baseUrl}/api?api_id=e7b2e813`, dbData, {
      headers,
    });
    if (data.code === 0) {
      commit('setDBs', [...state.dbs, data.data]);
    }
    return data;
  },
  async removeDb({ commit, state, rootState, dispatch }, dbId) {
    const headers = await dispatch('getHeaders', null, { root: true });
    const { data } = await axios.post(`${rootState.baseUrl}/api?api_id=b2e81aca`, { db_id: dbId }, {
      headers,
    });
    if (data.code === 0) {
      commit('setDBs', state.dbs.filter(p => p.id !== dbId));
    }
    return data;
  },
  async updateDb({ commit, state, rootState, dispatch }, dbData) {
    const headers = await dispatch('getHeaders', null, { root: true });
    const { data } = await axios.post(`${rootState.baseUrl}/api?api_id=1e1c616e`, dbData, {
      headers,
    });
    if (data.code === 0) {
      commit('setDBs', [...state.dbs, data.data]);
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