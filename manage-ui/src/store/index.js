// src/store/index.js  
import { createStore } from 'vuex';  

import api from './modules/api';  
import project from './modules/project';  
import db from './modules/db';  
import authkey from './modules/authkey';  

export default createStore({  
  state: {  
    baseUrl :'',
    authKey: localStorage.getItem("authkey"),
    Token: localStorage.getItem("token")
  },
  mutations: {  
  },  
  actions: {  
    getHeaders({ state }) {
      return {  
        'Content-Type': 'application/json',  
        authkey: state.authKey,  
        token: state.huhangToken,  
      };  
    },
  },  
  modules: {
    api,
    project,
    authkey,
    db
  }, 
});  
