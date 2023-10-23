// src/store/modules/api.js  
import axios from 'axios';  
  
const state = {  
  apis: [], // API列表  
};  
  
const mutations = {  
  setApis(state, apis) {  
    state.apis = apis;  
  },  
};  
  
const actions = {  
  async fetchAPIsPaged({ commit, dispatch, rootState }, { limit }) {  
    const headers = await dispatch('getHeaders', null, { root: true });  
    const { data } = await axios.post(  
      `${rootState.baseUrl}/api?api_id=8181e1ca`,  
      { limit },  
      { headers, }  
    );  
    commit('setApis', data.data);  
    return data;  
  },  
  async addAPI({ commit, state, rootState ,dispatch }, apiData) {  
    const headers = await dispatch('getHeaders', null, { root: true });  
    const { data } = await axios.post(`${rootState.baseUrl}/api?api_id=7b2e8161`, apiData, {  
      headers,  
    });  
    if (data.code === 0) {  
      commit('setApis', [...state.apis, data.data]);  
    }  
    return data;  
  },  
  async removeAPI({ commit, state, rootState,dispatch }, apiId) {  
    const headers = await dispatch('getHeaders', null, { root: true });  
    const { data } = await axios.post(`${rootState.baseUrl}/api?api_id=ab2e81aa`, { api_id: apiId }, {  
      headers,  
    });  
    if (data.code === 0) {  
      commit('setApis', state.apis.filter(p => p.id !== apiId));  
    }  
    return data;  
  },  
  async updateAPI({ commit, state, rootState, dispatch }, apiData) {  
    const headers = await dispatch('getHeaders', null, { root: true });  
    const { data } = await axios.post(`${rootState.baseUrl}/api?api_id=1616e7b2`, apiData, {  
      headers,  
    });  
    if (data.code === 0) {  
      commit('setApis', [...state.apis, data.data]);  
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
