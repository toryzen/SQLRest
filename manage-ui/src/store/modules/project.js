// src/store/modules/project.js
import axios from 'axios';

const state = {
  projects: [], // 项目列表
};

const mutations = {
  setProjects(state, projects) {
    state.projects = projects;
  },
};

const actions = {
  async fetchProjectsPaged({ commit, dispatch, rootState }, { limit }) {
    const headers = await dispatch('getHeaders', null, { root: true });
    const { data } = await axios.post(
      `${rootState.baseUrl}/api?api_id=c616e72a`,
      { limit },
      { headers }
    );
    commit('setProjects', data.data);
    return data;
  },
  async addProject({ commit, state, rootState, dispatch }, projectData) {
    const headers = await dispatch('getHeaders', null, { root: true });
    const { data } = await axios.post(`${rootState.baseUrl}/api?api_id=611c6164`, projectData, {
      headers,
    });
    if (data.code === 0) {
      commit('setProjects', [...state.projects, data.data]);
    }
    return data;
  },
  async removeProject({ commit, state, rootState, dispatch }, projectId) {
    const headers = await dispatch('getHeaders', null, { root: true });
    const { data } = await axios.post(`${rootState.baseUrl}/api?api_id=e8181ee7`, { project_id: projectId }, {
      headers,
    });
    if (data.code === 0) {
      commit('setProjects', state.projects.filter(p => p.id !== projectId));
    }
    return data;
  },
  async updateProject({ commit, state, rootState, dispatch }, projectData) {
    const headers = await dispatch('getHeaders', null, { root: true });
    const { data } = await axios.post(`${rootState.baseUrl}/api?api_id=e7b2e812`, projectData, {
      headers,
    });
    if (data.code === 0) {
      commit('setProjects', state.projects.map(p => (p.id === projectData.id ? data.data : p)));
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