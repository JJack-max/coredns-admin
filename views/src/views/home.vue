<template>
  <el-container>
    <el-header>
      <div class="container">
        <h1>CoreDNS</h1>
      </div>
    </el-header>
    <el-container style="  margin-top: 24px;">
      <el-aside width="200px">
        <el-menu default-active="/records">
          <el-submenu index="/">
            <template slot="title">
              <i class="el-icon-eleme"></i>
              <span slot="title">etcd</span>
            </template>

            <el-menu-item index="/records">记录</el-menu-item>
            <el-menu-item index="1-2">域名</el-menu-item>
          </el-submenu>
        </el-menu>
      </el-aside>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
// @ is an alias to /src

export default {
  name: 'Home'
}
</script>

<style lang="scss" scoped>
.container {
  height: 100%;

  box-sizing: border-box;
  border-bottom: 1px solid #dcdfe6;
  h1 {
    line-height: 60px;
    margin: 0;
    float: left;
    color: #409eff;
    font-size: 30px;
    font-weight: 500;
  }
}
.el-menu {
  height: 100%;
}
</style>
