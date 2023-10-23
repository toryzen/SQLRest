<!-- views/ProjectList.vue -->  
<template>  
  <div>  
    <h2>项目管理</h2>
 
    <el-row type="flex" justify="center" gutter="10">  
      <el-col :span="2">  
        <el-input v-model="searchParams.id" placeholder="ProjectID"></el-input>  
      </el-col>  
      <el-col :span="4">  
        <el-input v-model="searchParams.project_name" placeholder="项目名称"></el-input>  
      </el-col>  
      <el-col :span="4">  
        <el-input v-model="searchParams.memo" placeholder="项目备注"></el-input>  
      </el-col>  
      <el-col :span="2">  
        <el-button type="primary" @click="searchProjects" >搜索</el-button>  
      </el-col>  
      <el-col :span="1">  
        <el-button type="success" @click="showAddProjectForm">添加项目</el-button>  
      </el-col>  
    </el-row>  

  
    <!-- 添加项目表单 -->  
    <el-dialog title="添加项目" v-model="addFormVisible">  
      <el-form @submit.prevent="addProject" class="add-project-form">  
        <el-form-item label="项目名称">  
          <el-input v-model="newProject.project_name" placeholder="项目名称" required></el-input>  
        </el-form-item>  
        <el-form-item label="项目备注">  
          <el-input v-model="newProject.memo" placeholder="项目备注"></el-input>  
        </el-form-item>  
        <el-button type="primary" @click="addProject">添加</el-button>  
        <el-button @click="cancelAdd">取消</el-button>  
      </el-form>  
    </el-dialog>  
  
    <el-divider />  

    <el-table :data="projects">  
      <el-table-column type="index" label="ID"></el-table-column>  
      <el-table-column prop="project_id" label="ProjectID"></el-table-column> 
      <el-table-column prop="project_name" label="项目名称"></el-table-column>  
      <el-table-column prop="memo" label="备注"></el-table-column>  
      <el-table-column label="API Docs">
        <template #default="{ row }">  
          <a :href="getDocsUrl(row.project_id)" target="_blank">查看</a>  
        </template> 
      </el-table-column>  
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
          <el-button type="warning" @click="showUpdateForm(scope.row.id)"  size="small">更新</el-button>  
          <el-button type="danger" @click="confirmDelete(scope.row.id)"  size="small">删除</el-button>  
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
  
    <el-dialog title="更新项目信息" v-model="updateFormVisible" >  
      <el-form @submit.prevent="updateProjectInfo" class="update-project-form">  
        <el-form-item label="名称">
          <el-input v-model="updatedProject.project_name" placeholder="项目名称" required />  
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="updatedProject.memo" placeholder="备注" /> 
        </el-form-item>
        <el-input type="hidden" v-model="updatedProject.id"></el-input>  
        <el-button type="primary" @click="updateProjectInfo">保存</el-button>  
        <el-button @click="cancelUpdate">取消</el-button>  
      </el-form>  
    </el-dialog>  
    
  </div>  
</template>  
  
<script>  
export default {  
  name: "ProjectList",  
  data() {  
    return {  
      searchParams:{
        project_name: "",  
        id: "", 
        memo:""
      },
      newProject: {  
        project_name: "",  
        memo: "",  
      },  
      updateFormVisible: false,  
      updatedProject: {  
        id: "",  
        project_name: "",  
        memo: "",  
      },  
      currentPage: 1,  
      totalPages: 0,
      pageSize:10,
      totalRecords:0,
      addFormVisible: false, 
    };  
  },  
  computed: {  
    projects() {  
      return this.$store.state.project.projects;  
    },  
  },  
  methods: {  
    async addProject() {  
      if (this.newProject.project_name.trim() === "") { 
        this.$message({ type: "error", message: "项目名称不能为空", }); 
        return; 
      }
      const response = await this.$store.dispatch("project/addProject", this.newProject); 
      if (response.code === 200) {
        console.log(response)
        this.$message({ type: "success", message: "项目添加成功", });
      } else {
        this.$message({ type: "error", message: "项目添加失败", });
      }
      this.newProject = { project_name: "", memo: "" };  
      this.fetchProjectListByPage(); 
      this.addFormVisible = false;  
    },  
    async removeProject(projectId) {  
      const response = await this.$store.dispatch("project/removeProject", projectId); 
      if (response.code === 200) {
        this.$message({ type: "success", message: "项目删除成功", }); 
      } else {
        this.$message({ type: "error", message: "项目删除失败", });
      }
      this.fetchProjectListByPage();  
    },  
    confirmDelete(projectId) {  
      this.$confirm("确定要删除此项目吗？", "提示", {  
        confirmButtonText: "确定",  
       
        cancelButtonText: "取消",  
        type: "warning",  
      })  
        .then(() => {  
          this.removeProject(projectId);  
        })  
        .catch(() => {});  
    },  
    showUpdateForm(projectId) {  
      const projectIndex = this.projects.findIndex((p) => p.id === projectId);  
      if (projectIndex !== -1) {  
        const project = this.projects[projectIndex];  
        this.updatedProject = {  
          id: project.id,  
          project_name: project.project_name,  
          memo: project.memo,  
        };  
        this.updateFormVisible = true;  
      }       
    },
    async updateProjectInfo() {
      if (this.updatedProject.project_name.trim() === "") { 
        this.showMessage("error", "项目名称不能为空"); 
        return; 
      }
      const response = await this.$store.dispatch("project/updateProject", this.updatedProject); 
      if (response.code === 200) {
        this.$message({ type: "success", message: "项目更新成功", }); 
      } else {
        this.$message({ type: "error", message: "项目更新失败", });
      }
      this.updatedProject = { id: "", project_name: "", memo: "" };  
      this.updateFormVisible = false;  
      this.fetchProjectListByPage();  
    },  
    cancelUpdate() {  
      this.updateFormVisible = false;  
    },  
    async fetchProjectListByPage() {  
      const pageSize = this.pageSize;  
      const offset = (this.currentPage - 1) * pageSize;  
      let limit = "";  
      if (this.searchParams.id) {
        limit += ` AND project_id='${this.searchParams.id}'`;
      }
      if (this.searchParams.project_name) {
        limit += ` AND project_name LIKE '%${this.searchParams.project_name}%'`;
      }
      if (this.searchParams.memo) {
        limit += ` AND memo LIKE '%${this.searchParams.memo}%'`;
      }
      limit += ` limit ${offset},${pageSize}`;  

      try {  
        const data = await this.$store.dispatch("project/fetchProjectsPaged", {    
          limit,    
        });    
        this.totalRecords = data.data.length ? parseInt(data.data[0].allcnt) : 0;    
        this.totalPages = Math.ceil(this.totalRecords / pageSize);    
      } catch (error) {  
        if (error.response && error.response.status === 401) {  
          localStorage.setItem('authkey', '');  
          this.$message({message: '授权认证失败，请重试！', type: 'warning',showClose: true,duration: 0});
          setTimeout(() => { window.location.reload() }, 3000); 
        } else {  
          this.$message({message: error, type: 'error',showClose: true,duration: 0});
          setTimeout(() => {  window.location.reload() }, 3000); 
        }
      }
    },
    handleSizeChange(val) {  
      this.currentPage = 1
      this.pageSize = val;
      this.fetchProjectListByPage(); 
    },  
    handleCurrentChange(val) {  
      this.currentPage = val;  
      this.fetchProjectListByPage();  
    },  
    // 显示添加项目表单  
    showAddProjectForm() {  
        this.addFormVisible = true;  
    },  
    // 取消添加项目  
    cancelAdd() {  
      this.addFormVisible = false;  
    },
    async searchProjects() {
      this.currentPage = 1;
      this.fetchProjectListByPage();
    },
    getDocsUrl(id){
      return this.$store.state.baseUrl+"/help?project_id="+id
    },
  },  
  created() {  
    this.fetchProjectListByPage(this.currentPage);  
  },  
  
};  
</script>  

<style scoped>
.pagination-container { 
  margin-top: 20px; 
  display: flex; 
  justify-content: center; 
  }
  a{
  text-decoration: none;
  color: #2c3e50;
}

</style>