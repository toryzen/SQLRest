<template>
  <div>
    <h2>数据源管理</h2>

    <!-- 搜索表单 -->
    <el-row type="flex" justify="center" gutter="10">
      <el-col :span="4">
        <el-input v-model="searchParams.dbname" placeholder="数据库名"></el-input>
      </el-col>
      <el-col :span="4">
        <el-input v-model="searchParams.memo" placeholder="备注"></el-input>
      </el-col>
      <el-col :span="2">
        <el-button type="primary" @click="searchDBs">搜索</el-button>
      </el-col>
      <el-col :span="1">
        <el-button type="success" @click="showAddDBForm">添加数据库</el-button>
      </el-col>
    </el-row>
    
    <!-- 添加数据库表单 -->
    <el-dialog title="添加数据库" v-model="addFormVisible">
      <el-form @submit.prevent="addDB" class="add-db-form">
        <el-form-item label="数据库名">
          <el-input v-model="newDB.dbname" placeholder="数据库名" required></el-input>
        </el-form-item>
        <el-form-item label="IP">
          <el-input v-model="newDB.ip" placeholder="IP" required></el-input>
        </el-form-item>
        <el-form-item label="端口">
          <el-input v-model="newDB.port" placeholder="端口" required></el-input>
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="newDB.user" placeholder="用户名" required></el-input>
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="newDB.pwd" placeholder="密码" required></el-input>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="newDB.memo" placeholder="备注"></el-input>
        </el-form-item>
        <el-button type="primary" @click="addDB">添加</el-button>
        <el-button @click="cancelAdd">取消</el-button>
      </el-form>
    </el-dialog>

    <el-divider />

    <el-table :data="dbs">
      <el-table-column type="index" label="ID"></el-table-column>
      <el-table-column prop="dbname" label="数据库名"></el-table-column>
      <el-table-column prop="ip" label="IP"></el-table-column>
      <el-table-column prop="port" label="端口"></el-table-column>
      <el-table-column prop="user" label="用户名"></el-table-column>
      <el-table-column prop="pwd" label="密码"></el-table-column>
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

    <el-dialog title="更新数据库信息" v-model="updateFormVisible" >
      <el-form @submit.prevent="updateDBInfo" class="update-db-form">
        <el-form-item label="用户名">
          <el-input v-model="updatedDB.user" placeholder="用户名" required></el-input>
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="updatedDB.pwd" placeholder="密码" required ></el-input>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="updatedDB.memo" placeholder="备注"></el-input>
        </el-form-item>
        <el-form-item label="数据库名">
          <el-input v-model="updatedDB.dbname" placeholder="数据库名" required></el-input>
        </el-form-item>
        <el-form-item label="IP">
          <el-input v-model="updatedDB.ip" placeholder="IP" required></el-input>
        </el-form-item>
        <el-form-item label="端口">
          <el-input v-model="updatedDB.port" placeholder="端口" required></el-input>
        </el-form-item>
        <el-input type="hidden" v-model="updatedDB.id"></el-input>
        <el-button type="primary" @click="updateDBInfo">保存</el-button>
        <el-button @click="cancelUpdate">取消</el-button>
      </el-form>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'DBList',
  data() {
    return {
      searchParams:{
        memo: '',
        dbname: '',
      },
      newDB: {
        user: '',
        pwd: '',
        memo: '',
        dbname: '',
        ip: '',
        port: '',
      },
      updateFormVisible: false,
      updatedDB: {
        id: '',
        user: '',
        pwd: '',
        memo: '',
        dbname: '',
        ip: '',
        port: '',
      },
      currentPage: 1,
      totalPages: 0,
      pageSize: 10,
      totalRecords: 0,
      addFormVisible:false,
    };
  },
  computed: {
    dbs() {
      return this.$store.state.db.dbs;
    },
  },
  methods: {
    async addDB() {
      const response = await this.$store.dispatch('db/addDb', this.newDB);
      if (response.code === 200) {
        console.log(response)
        this.$message({ type: "success", message: "添加成功", });
      } else {
        this.$message({ type: "error", message: "添加失败", });
      }
      this.newDB = { user: '', pwd: '', memo: '', dbname: '', ip: '', port: '' };
      this.fetchDBListByPage();
      this.addFormVisible = false;  
    },
    async removeDB(dbId) {
      const response =await this.$store.dispatch('db/removeDb', dbId);
      if (response.code === 200) {
        console.log(response)
        this.$message({ type: "success", message: "删除成功", });
      } else {
        this.$message({ type: "error", message: "删除失败", });
      }
      this.fetchDBListByPage();
    },
    confirmDelete(dbId) {
      this.$confirm("确定要删除此数据库吗？", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(() => {
          this.removeDB(dbId);
        })
        .catch(() => {});
    },
    showUpdateForm(dbId) {
      const db = this.dbs.find((d) => d.id === dbId);
      this.updatedDB = { id: db.id, user: db.user, pwd: db.pwd, memo: db.memo, dbname: db.dbname, ip: db.ip, port: db.port };
      this.updateFormVisible = true;
    },
    async updateDBInfo() {
      const response =await this.$store.dispatch('db/updateDb', this.updatedDB);
      if (response.code === 200) {
        console.log(response)
        this.$message({ type: "success", message: "更新成功", });
      } else {
        this.$message({ type: "error", message: "更新失败", });
      }
      this.updatedDB = { id: '', user: '', pwd: '', memo: '', dbname: '', ip: '', port: '' };
      this.updateFormVisible = false;
      this.fetchDBListByPage();
    },
    cancelUpdate() {
      this.updateFormVisible = false;
    },
    async fetchDBListByPage() {
      const pageSize = this.pageSize;
      const offset = (this.currentPage - 1) * pageSize;
      let limit = "";  
      if (this.searchParams.dbname) {
        limit += ` AND dbname LIKE '%${this.searchParams.dbname}%'`;
      }
      if (this.searchParams.memo) {
        limit += ` AND memo LIKE '%${this.searchParams.memo}%'`;
      }
      limit += ` limit ${offset},${pageSize}`; 
      try{
        const data = await this.$store.dispatch("db/fetchDBsPaged", {
          limit,
        });
        this.totalRecords = data.data.length ? parseInt(data.data[0].allcnt) : 0;
        this.totalPages = Math.ceil(this.totalRecords / pageSize);
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
      this.fetchDBListByPage();
    },
    handleCurrentChange(val) {
      this.currentPage = val;
      this.fetchDBListByPage();
    },
    // 显示添加项目表单  
    showAddDBForm() {  
      this.addFormVisible = true;  
    },  
    // 取消添加项目  
    cancelAdd() {  
      this.addFormVisible = false;  
    },
    async searchDBs() {
      this.currentPage = 1;
      this.fetchDBListByPage();
    },
  },
  created() {
    this.fetchDBListByPage();
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

.add-db-form {
  margin-bottom: 20px;
}
</style>