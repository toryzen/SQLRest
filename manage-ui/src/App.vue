<template>

<div class="common-layout">
    <el-container>
      <el-header>
        <el-menu  class="el-menu-demo" mode="horizontal" :ellipsis="false" >
          <el-menu-item index="0"><router-link to="/">SQLRest Manage UI</router-link></el-menu-item>
          
          <div class="flex-grow" />
           
            <el-menu-item v-if="!authKey" index="6"  @click="showLoginDialog">登录</el-menu-item>

            <el-menu-item  v-if="authKey" index="1"><router-link to="/project-list">项目</router-link></el-menu-item>
            <el-menu-item  v-if="authKey" index="2"><router-link to="/api-list">API</router-link></el-menu-item>
            <el-menu-item  v-if="authKey" index="3"><router-link to="/authkey-list">密钥</router-link></el-menu-item>
            <el-menu-item  v-if="authKey" index="4"><router-link to="/db-list">数据源</router-link></el-menu-item>
            <el-sub-menu  v-if="authKey" index="5">
              <template #title>操作</template>
              <el-menu-item index="5-2"  @click="logout"> 退出AuthKey：{{ authKey }} </el-menu-item>
            </el-sub-menu>
        </el-menu>
      </el-header>
      <el-main>
        <div v-if="!authKey">          
          <h2>Project API Search</h2>
          <el-row type="flex" justify="center" gutter="10">
            <el-col :span="5">
              <el-input v-model="inputProjectID" placeholder="请输入 Project ID"></el-input>
            </el-col>
            <el-col :span="2">
              <el-button type="primary" @click="searchProject">搜索</el-button>
            </el-col>
          </el-row>
        </div>
        <div v-else>
          <router-view/>
        </div>
      </el-main>
      <!--
      <el-footer>
        <div v-if="!authKey">
          <el-divider />
          <el-text class="mx-1"><div class="nav"><router-link to="/readme"  target="_blank"> - 帮助手册 - </router-link></div></el-text>
        </div>
      </el-footer>
      -->

      <!-- 登录对话框 -->
      <el-dialog title="请输入管理密钥" v-model="loginDialogVisible">
        <el-row type="flex" justify="center" gutter="10">
          <el-col :span="8">
            <el-input v-model="inputAuthKey" placeholder="请输入AuthKey" show-password></el-input>
          </el-col>
          <el-col :span="2">
            <el-button type="primary" @click="saveAuthKey">登录</el-button>
          </el-col>
        </el-row>
        <br>
      </el-dialog>

    </el-container>
  </div>
</template>

<script>
export default {
  data() {
    return {
      authKey: this.$store.state.authKey,
      inputAuthKey: '',
      inputProjectID: '',
      loginDialogVisible: false,
    };
  },
  methods: {
    showLoginDialog() {
      this.loginDialogVisible = true;
    },
    async saveAuthKey() {
      localStorage.setItem('authkey', this.inputAuthKey);
      this.$store.state.authKey = this.inputAuthKey
      let limit = "limit 1";  
      try {  
        await this.$store.dispatch("project/fetchProjectsPaged", {    
          limit,    
        });    
        location.reload(); 
      } catch (error) {  
        if (error.response && error.response.status === 401) {  
          localStorage.setItem('authkey', '');  
          this.$message({message: '授权认证失败，请重试！', type: 'warning',showClose: true});
          //setTimeout(() => { window.location.reload() }, 3000); 
        } else {  
          this.$message({message: error, type: 'error',showClose: true});
          //setTimeout(() => {  window.location.reload() }, 3000); 
        }
      }
    },
    logout() {  
      localStorage.removeItem('authkey');  
      location.reload();  
    },
    async searchProject() {  
      try {  
        const response = await fetch(`${this.$store.state.baseUrl}/help?project_id=${this.inputProjectID}&check=check`);  
        const result = await response.text();  
      
        if (result === 'ok') {  
          window.open(`${this.$store.state.baseUrl}/help?project_id=${this.inputProjectID}`, '_blank');  
        } else {  
          this.$message.error('未找到该项目，请检查输入的Project ID是否正确。');  
        }  
      } catch (error) {  
        this.$message.error(`请求发生错误: ${error}`);  
      }  
    } 
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

a{
  text-decoration: none;
}


.nav a {
  font-weight: bold;
  color: #2c3e50;
}

.nav a.router-link-exact-active {
  color: #42b983;
}

.flex-grow {
  flex-grow: 1;
}
</style>