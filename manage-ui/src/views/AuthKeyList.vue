<!-- views/AuthKeyList.vue -->  
<template>
  <div>
    <h2>密钥管理</h2>

    <!-- 搜索部分 -->  
    <el-row type="flex" justify="center" gutter="10">  
      <el-col :span="3">  
        <el-input v-model="searchParams.authkey" placeholder="AuthKey"></el-input>  
      </el-col>  
      <el-col :span="3">  
        <el-select v-model="searchParams.project_id" placeholder="项目" filterable>  
          <el-option  
            v-for="item in projectList"  
            :key="item.project_id"  
            :label="item.project_name"  
            :value="item.project_id">  
          </el-option>  
        </el-select>  
      </el-col>  
      <el-col :span="4">  
        <el-input v-model="searchParams.memo" placeholder="备注"></el-input>  
      </el-col>  
      <el-col :span="2">  
        <el-button type="primary" @click="searchAuthKeys">搜索</el-button>  
      </el-col>  
      <el-col :span="2">  
        <el-button type="success" @click="showAddAuthkeyForm()">添加AuthKey</el-button>  
      </el-col>  
    </el-row>  
  
    <!-- 添加 AuthKey 对话框 -->  
    <el-dialog title="添加 AuthKey" v-model="addFormVisible">  
      <el-form @submit.prevent="addAuthKey" class="add-authkey-form">  
        <el-form-item label="项目">  
          <el-select v-model="newAuthKey.project_id" placeholder="项目ID" required filterable @change="fetchAPIList">  
            <el-option  
              v-for="item in projectList"  
              :key="item.project_id"  
              :label="item.project_name"  
              :value="item.project_id">  
            </el-option>  
          </el-select>  
        </el-form-item>  
        <el-form-item label="API IDs">  
          <el-select v-model="newAuthKey.api_ids" placeholder="API IDs" multiple filterable>  
            <el-option v-if="newAuthKey.project_id" :label="'全部'" :value="'-1'"></el-option>  
            <el-option  
              v-for="item in apiList"  
              :key="item.api_id"  
              :label="item.apiname"  
              :value="item.api_id">  
            </el-option>  
          </el-select>  
        </el-form-item>  
        <el-form-item label="备注">  
          <el-input v-model="newAuthKey.memo" placeholder="备注"></el-input>  
        </el-form-item>  
        <el-button type="primary" @click="addAuthKey">添加</el-button>  
        <el-button @click="cancelAdd">取消</el-button>  
      </el-form>  
    </el-dialog>  

    <el-divider />

    <el-table :data="authKeys">
      <el-table-column type="index" label="ID"></el-table-column>
      <el-table-column prop="authkey" label="AuthKey"></el-table-column>
      <el-table-column prop="project_name" label="项目"></el-table-column>
      <el-table-column prop="api_ids" label="API IDs"></el-table-column>
      <el-table-column prop="memo" label="备注"></el-table-column>
      <el-table-column label="信息">
        <template #default="{ row }">
          <el-tooltip placement="top">
            <template #content>
              最后更新人：{{ row.modified_user }} <br>
              更新时间：{{ row.modified_stime }} <br>
              创建人：{{ row.created_user }} <br>
              创建时间：{{ row.created_stime }} <br>
            </template>
            <el-button  type="info" size="small">信息</el-button>  
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <template v-slot:default="scope">
          <div class="action-buttons">  
            <el-button type="warning" @click="showUpdateForm(scope.row.authkey)"  size="small">更新</el-button>
            <el-button type="danger" @click="confirmDelete(scope.row.authkey)"  size="small">删除</el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <el-divider />
    <div class="pagination-container">
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="currentPage"
        :page-sizes="[10, 20, 50, 100]"
        :page-size="pageSize"
        layout="total, sizes, prev, pager, next, jumper"
        :total="totalRecords"
        >
      </el-pagination>
    </div>

    <el-dialog title="更新AuthKey 信息" v-model="updateFormVisible" >
      <el-form @submit.prevent="updateAuthKeyInfo" class="update-authkey-form">
        <el-form-item label="AuthKey">
          <el-input v-model="updatedAuthKey.authkey" placeholder="AuthKey" readonly required></el-input>
        </el-form-item>
        <el-form-item label="API IDs">
          <!--<el-input v-model="updatedAuthKey.api_ids" placeholder="API IDs" ></el-input>-->
          <el-select v-model="updatedAuthKey.api_ids" placeholder="API IDs" multiple filterable>  
            <el-option :label="'全部'" :value="'-1'"></el-option> 
            <el-option  
              v-for="item in apiList"  
              :key="item.api_id"  
              :label="item.apiname"  
              :value="item.api_id">  
            </el-option>  
          </el-select>  
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="updatedAuthKey.memo" placeholder="备注"></el-input>
        </el-form-item>
        <el-input type="hidden" v-model="updatedAuthKey.id"></el-input>
        <el-button type="primary" @click="updateAuthKeyInfo">保存</el-button>
        <el-button @click="cancelUpdate">取消</el-button>
      </el-form>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'AuthKeyList',
  data() {
    return {
      apiList:[],
      projectList: [],
      searchParams:{
        project_id: '',
        authkey: '',
        memo: '',
      },
      newAuthKey: {
        project_id: '',
        api_ids: '',
        memo: '',
      },
      updateFormVisible: false,
      updatedAuthKey: {
        id: '',
        authkey: '',
        project_id: '',
        api_ids: '',
        memo: '',
      },
      currentPage: 1,
      totalPages: 0,
      pageSize:10,
      totalRecords:0,
      addFormVisible:false,
    };
  },
  computed: {
    authKeys() {
      return this.$store.state.authkey.authKeys;
    },
  },
  methods: {
    async addAuthKey() {
      this.newAuthKey.api_ids = this.newAuthKey.api_ids.join(',');
      const response = await this.$store.dispatch('authkey/addAuthKey', this.newAuthKey);
      if (response.code === 200) {
        console.log(response)
        this.$message({ type: "success", message: "添加成功", });
      } else {
        this.$message({ type: "error", message: "添加失败", });
      }
      this.newAuthKey = { project_id: '', api_ids: '', memo: '' };
      this.fetchAuthKeyListByPage();
      this.addFormVisible = false; 
    },
    async removeAuthKey(authKeyId) {
      const response = await this.$store.dispatch('authkey/removeAuthKey', authKeyId);
      if (response.code === 200) {
        console.log(response)
        this.$message({ type: "success", message: "删除成功", });
      } else {
        this.$message({ type: "error", message: "删除失败", });
      }
      this.fetchAuthKeyListByPage();
      this.addFormVisible = false;  
    },
    confirmDelete(authKeyId) {
      this.$confirm("确定要删除此 AuthKey 吗？", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(() => {
          this.removeAuthKey(authKeyId);
        })
        .catch(() => {});
    },
    showUpdateForm(authKeyId) {
      const authKey = this.authKeys.find((a) => a.authkey === authKeyId);  
      this.fetchAPIList(authKey.project_id);  
      const api_ids = authKey.api_ids.split(',').map(String);  
      this.updatedAuthKey = {  
        id: authKey.id,  
        authkey: authKey.authkey,  
        api_ids: api_ids,  
        memo: authKey.memo,  
      };  
      this.updateFormVisible = true; 
    },
    async updateAuthKeyInfo() {
      this.updatedAuthKey.api_ids = this.updatedAuthKey.api_ids.join(',');  
      const response = await this.$store.dispatch('authkey/updateAuthKey', this.updatedAuthKey);
      console.log(response)
      if (response.code === 200) {
        this.$message({ type: "success", message: "更新成功", });
      } else {
        this.$message({ type: "error", message: "更新失败", });
      }
      this.updatedAuthKey = { id: '', authkey: '', api_ids: '', memo: '' };
      this.updateFormVisible = false;
      this.fetchAuthKeyListByPage();
    },
    cancelUpdate() {
      this.updateFormVisible = false;
    },
    async fetchAuthKeyListByPage() {
      const pageSize = this.pageSize;
      const offset = (this.currentPage - 1) * pageSize;
      let limit = "";  
      if (this.searchParams.project_id) {
        limit += ` AND project_id= '${this.searchParams.project_id}' `;
      }
      if (this.searchParams.authkey) {
        limit += ` AND authkey = '${this.searchParams.authkey}'`;
      }
      if (this.searchParams.memo) {
        limit += ` AND memo LIKE '%${this.searchParams.memo}%'`;
      }
      limit += ` limit ${offset},${pageSize}`;
      try{
        const data = await this.$store.dispatch("authkey/fetchAuthKeysPaged", {
          limit,
        });
        this.totalRecords = data.data.length ? parseInt(data.data[0].allcnt) : 0;
        this.totalPages = Math.ceil(this.totalRecords / pageSize);
        this.fetchProjectList()
      }catch (error) { 
        if (error.response && error.response.status === 401) {  
          localStorage.setItem('authkey', '');  
          this.$message({message: '授权认证失败，请重试！', type: 'warning',showClose: true,duration: 0});
          setTimeout(() => { window.location.reload() }, 3000); 
        } else {  
          this.$message({message: error, type: 'error',showClose: true,duration: 0});
          setTimeout(() => { window.location.reload() }, 3000); 
        }
      } 
      
    },
    handleSizeChange(val) {
      this.currentPage = 1
      this.pageSize = val;
      this.fetchAuthKeyListByPage();
    },
    handleCurrentChange(val) {
      this.currentPage = val;
      this.fetchAuthKeyListByPage();
    },
    async fetchProjectList(){
      const limit = "";  
      const data = await this.$store.dispatch("project/fetchProjectsPaged", {  
        limit,  
      });
      this.projectList = data.data;
    },
    async fetchAPIList(pid) {  
      const project_id = pid?pid:this.newAuthKey.project_id;  
      const limit = project_id ? ` AND project_id='${project_id}'` : "";    
      const data = await this.$store.dispatch("api/fetchAPIsPaged", {  
        limit,
      });  
      this.apiList = data.data;
    },  
    showAddAuthkeyForm() {  
      this.addFormVisible = true;  
    },  
    cancelAdd() {  
      this.addFormVisible = false;  
    },
    async searchAuthKeys() {
      this.currentPage = 1;
      this.fetchAuthKeyListByPage();
    },
  },
  created() {
    this.fetchAuthKeyListByPage(this.currentPage);
  },
};
</script>

<style scoped>
.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.action-buttons {  
  display: flex;  
  flex-direction: row;  
  align-items: center;  
}  

</style>
