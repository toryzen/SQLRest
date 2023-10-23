<!-- views/APIList.vue -->  
<template>
  <div>
    <h2>API 管理</h2>

    <!-- 搜索栏 -->
    <el-row type="flex" justify="center" gutter="10">
      <el-col :span="2">
        <el-input v-model="searchParams.api_id" placeholder="ApiId"></el-input>
      </el-col>
      <el-col :span="3">
        <el-input v-model="searchParams.apiname" placeholder="API名称"></el-input>
      </el-col>
      <el-col :span="3">
        <el-select v-model="searchParams.project_id" placeholder="项目" filterable clearable>
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
        <el-button type="primary" @click="searchAPIs">搜索</el-button>
      </el-col>
      <el-col :span="1">
        <el-button type="success" @click="showAddAPIForm()">添加API</el-button>
      </el-col>
    </el-row>

    <!-- 添加API表单 -->
    <el-dialog title="添加API" v-model="addFormVisible">
      <el-form @submit.prevent="addAPI" class="add-api-form">
        <el-form-item label="API名称">
          <el-input v-model="newAPI.apiname" placeholder="API名称" required></el-input>
        </el-form-item>
        <el-form-item label="项目">
          <el-select v-model="newAPI.project_id" placeholder="项目" required filterable>
            <el-option
              v-for="item in projectList"
              :key="item.project_id"
              :label="item.project_name"
              :value="item.project_id">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="数据库ID">
          <el-select v-model="newAPI.db_id" placeholder="数据库" required filterable>
            <el-option
              v-for="item in dbList"
              :key="item.db_id"
              :label="item.dbname"
              :value="item.db_id">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="联接条件">
          <el-input v-model="newAPI.joint" placeholder="联接条件"></el-input>
        </el-form-item>
        <el-form-item label="源SQL">
          <el-input v-model="newAPI.sourcesql" placeholder="源SQL(请勿使用单引号)" required></el-input>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="newAPI.memo" placeholder="备注"></el-input>
        </el-form-item>
        <el-button type="primary" @click="addAPI">添加</el-button>
        <el-button @click="cancelAdd">取消</el-button>
      </el-form>
    </el-dialog>

    <el-divider />

    <el-table :data="apis">
      <el-table-column type="index" label="ID"></el-table-column>
      <el-table-column prop="api_id" label="ApiId"></el-table-column>
      <el-table-column prop="apiname" label="API名称"></el-table-column>
      <el-table-column prop="project_name" label="项目"></el-table-column>
      <el-table-column prop="dbname" label="数据源"></el-table-column>
      <el-table-column prop="joint" label="联合鉴权"></el-table-column>
      <el-table-column label="源SQL">
        <template v-slot:default="scope">
          <el-popover
              :width="400"
              placement="top"
              trigger="click"
              :content="scope.row.sourcesql"
            >
            <template v-slot:reference>
              <el-button size="default">查看源SQL</el-button>
            </template>
          </el-popover>
        </template>
      </el-table-column>
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
            <el-button type="warning" @click="showUpdateForm(scope.row.id)"  size="small">更新</el-button>
            <el-button type="danger" @click="confirmDelete(scope.row.id)"  size="small">删除</el-button>
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

    <el-dialog title="更新 API 信息" v-model="updateFormVisible">
      <el-form @submit.prevent="updateAPIInfo" class="update-api-form">
        <el-form-item label="API名称">
          <el-input v-model="updatedAPI.apiname" placeholder="API名称" required></el-input>
        </el-form-item>
        <el-form-item label="数据库ID">
          <el-select v-model="updatedAPI.db_id" placeholder="数据库" required filterable>
            <el-option
              v-for="item in dbList"
              :key="item.db_id"
              :label="item.dbname"
              :value="item.db_id">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="联接条件">
          <el-input v-model="updatedAPI.joint" placeholder="联接条件"></el-input>
        </el-form-item>
        <el-form-item label="源SQL">
          <el-input v-model="updatedAPI.sourcesql" placeholder="源SQL" type="textarea" required></el-input>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="updatedAPI.memo" placeholder="备注"></el-input>
        </el-form-item>
        <el-input type="hidden" v-model="updatedAPI.id"></el-input>
        <el-form-item>
          <el-button type="primary" @click="updateAPIInfo">保存</el-button>
          <el-button @click="cancelUpdate">取消</el-button>
        </el-form-item>
      </el-form>
    </el-dialog>
  </div>
</template>
  
<script>  
export default {  
  name: "APIList",  
  data() {  
    return {
      dbList: [],
      projectList: [],
      searchParams:{
        api_id: "",  
        apiname: "",  
        project_id: "",  
        memo: "",  
      },
      newAPI: {  
        apiname: "",  
        project_id: "",  
        db_id: "",  
        joint: "",  
        sourcesql: "",  
        memo: "",  
      },  
      updateFormVisible: false,  
      updatedAPI: {  
        id: "",  
        apiname: "",  
        project_id: "",  
        db_id: "",  
        joint: "",  
        sourcesql: "",  
        memo: "",  
      },  
      currentPage: 1,  
      totalPages: 0,  
      pageSize: 10,  
      totalRecords:0,
      addFormVisible:false,
    };  
  },  
  computed: {  
    apis() {  
      return this.$store.state.api.apis;  
    },  
  },  
  methods: {  
    async addAPI() {
      if (this.newAPI.apiname.trim() === "") {  
        this.$message({ type: "error", message: "API名称不能为空" });  
        return;  
      }  
      const response = await this.$store.dispatch("api/addAPI", this.newAPI);  
      if (response.code === 200) {
        console.log(response)
        this.$message({ type: "success", message: "添加成功", });
      } else {
        this.$message({ type: "error", message: "添加失败", });
      }
      this.newAPI = {  
        apiname: "",  
        project_id: "",  
        db_id: "",  
        joint: "",  
        sourcesql: "",  
        memo: "",  
      };  
      this.fetchAPIListByPage();
      this.addFormVisible = false; 
    },  
    async removeAPI(apiId) {  
      const response = await this.$store.dispatch("api/removeAPI", apiId);  
      if (response.code === 200) {
        console.log(response)
        this.$message({ type: "success", message: "删除成功", });
      } else {
        this.$message({ type: "error", message: "删除失败", });
      }
      this.fetchAPIListByPage();  
    },  
    confirmDelete(apiId) {  
      this.$confirm("确定要删除此 API 吗？", "提示", {  
        confirmButtonText: "确定",  
        cancelButtonText: "取消",  
        type: "warning",  
      })  
        .then(() => {  
          this.removeAPI(apiId);  
        })  
        .catch(() => {});  
    },  
    showUpdateForm(apiId) {  
      const api = this.apis.find((a) => a.id === apiId);  
      this.updatedAPI = {  
        id: api.id,  
        apiname: api.apiname,  
        project_id: api.project_id,  
        db_id: api.db_id,  
        joint: api.joint,  
        sourcesql: api.sourcesql,
        memo: api.memo,  
      };  
      this.updateFormVisible = true;  
    },  
    async updateAPIInfo() {  
      if (this.updatedAPI.apiname.trim() === "") {  
        this.$message({ type: "error", message: "API名称不能为空" });  
        return;  
      }  
      const response = await this.$store.dispatch("api/updateAPI", this.updatedAPI);  
      if (response.code === 200) {
        console.log(response)
        this.$message({ type: "success", message: "更新成功", });
      } else {
        this.$message({ type: "error", message: "更新失败", });
      }
      this.updatedAPI = {  
        id: "",  
        apiname: "",  
        project_id: "",  
        db_id: "",          joint: "",  
        sourcesql: "",  
        memo: "",  
      };  
      this.updateFormVisible = false;  
      this.fetchAPIListByPage();  
    },  
    cancelUpdate() {  
      this.updateFormVisible = false;  
    },  
    async fetchAPIListByPage() {  
      const pageSize = this.pageSize;  
      const offset = (this.currentPage - 1) * pageSize;  
      let limit = "";  
      if (this.searchParams.api_id) {
        limit += ` AND api_id='${this.searchParams.api_id}'`;
      }
      if (this.searchParams.project_id) {
        limit += ` AND project_id='${this.searchParams.project_id}'`;
      }
      if (this.searchParams.apiname) {
        limit += ` AND apiname LIKE '%${this.searchParams.apiname}%'`;
      }
      if (this.searchParams.memo) {
        limit += ` AND memo LIKE '%${this.searchParams.memo}%'`;
      }
      limit += ` limit ${offset},${pageSize}`;  
      try{
        const data = await this.$store.dispatch("api/fetchAPIsPaged", {  
          limit,
        });  
        this.totalRecords = data.data.length ? parseInt(data.data[0].allcnt) : 0;  
        this.totalPages = Math.ceil(this.totalRecords / pageSize);  
        this.fetchProjectList()
        this.fetchDBList()
      } catch (error) {  
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
      this.currentPage = 1;  
      this.pageSize = val;  
      this.fetchAPIListByPage();  
    },  
    handleCurrentChange(val) {  
      this.currentPage = val;  
      this.fetchAPIListByPage();  
    },
    async fetchProjectList(){
      const limit = "";  
      const data = await this.$store.dispatch("project/fetchProjectsPaged", {  
        limit,  
      });
      this.projectList = data.data;
    },
    async fetchDBList() {
      const limit = "";  
      const data = await this.$store.dispatch("db/fetchDBsPaged", {
        limit,
      });
      this.dbList = data.data;
    },
    // 显示添加项目表单  
    showAddAPIForm() {  
        this.addFormVisible = true;  
      },  
      // 取消添加项目  
      cancelAdd() {  
        this.addFormVisible = false;  
      },
      async searchAPIs() {
        this.currentPage = 1;
        this.fetchAPIListByPage();
      },
  },  
  created() {  
    this.fetchAPIListByPage(this.currentPage);  
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


       

