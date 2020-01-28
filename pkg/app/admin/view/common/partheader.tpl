{{define "common/partheader.tpl"}}
<header id="topbar" class="navbar navbar-default bs-docs-nav" role="banner">
  <div class="container">
    <div class="row">
    <div class="navbar-header">
      <button class="navbar-toggle collapsed" type="button" data-toggle="collapse" data-target=".bs-navbar-collapse">
        <span class="sr-only">导航</span>
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
      </button>
      <a style="font-size: 14px;" class="navbar-brand" rel="home" href="/" >
        <strong>{{.appname}}</strong>
      </a>
    </div>

    <nav class="collapse navbar-collapse bs-navbar-collapse" role="navigation" >
      <!-- 左侧主要菜单 -->
      <ul class="nav navbar-nav">
        <li>
          <a href='/'>
            <span class="glyphicon glyphicon-home"></span> 云盘管理后台首页
          </a>
        </li>
        <li>
          <a href="javascript:;" class="dropdown-toggle" data-hover="dropdown">
            <span class="glyphicon glyphicon-th-large"></span> 节点管理 <b class="caret"></b>
          </a>
          <ul class="dropdown-menu">
            <li role="presentation" class="dropdown-header">main route</li>
            <li><a href="/node/list">节点列表</a></li>
            <li role="presentation" class="divider"></li>
            <li role="presentation" class="dropdown-header">selected picsets</li>
            <li><a href="/file/list">文件管理</a></li>
          </ul>
        </li>
        <li>
          <a href="javascript:;" class="dropdown-toggle" data-hover="dropdown">
            <span class="glyphicon glyphicon-list-alt"></span> 会员管理 <b class="caret"></b>
          </a>
          <ul class="dropdown-menu">
            <li><a href="/user/list">会员列表</a></li>
            <li><a href="/userlevel/list">等级管理</a></li>
            <li><a href="/order/list">订单管理</a></li>
            <li><a href="/share/list">分享管理</a></li>
          </ul>
        </li>
        <li>
          <a href="javascript:;" class="dropdown-toggle" data-hover="dropdown">
            <span class="glyphicon glyphicon-cog"></span> 系统管理 <b class="caret"></b>
          </a>
          <ul class="dropdown-menu">
            <li><a href="/manager/list">管理员账号</a></li>
            <li><a href="/role/list">角色及权限管理</a></li>
          </ul>
        </li>
      </ul>

      <!-- 右侧菜单 -->
      <ul class="nav navbar-nav navbar-right">
        <li class="dropdown">
          <a href="javascript:;" role="button" class="dropdown-toggle" data-hover="dropdown">
            <span class='glyphicon glyphicon-user'></span> 
            {{if .IsLogin}}
              {{.LoginInfo.MgrEnt.Email}} 
            {{else}}
              账户
            {{end}}
            <b class="caret"></b>
          </a>
          <ul class="dropdown-menu">
            {{if .IsLogin}}
                <li ><a href='/changePwd'>
                  <span class='glyphicon glyphicon-certificate'></span> 修改密码
                </a></li>
                <li ><a href='/logout'>
                  <span class='glyphicon glyphicon-log-out'></span> 退出登录
                </a></li>
            {{else}}
                <li ><a href='/login'>
                  <span class='glyphicon glyphicon-log-in'></span> 登录
                </a></li>
                <!--
                <li ><a href='/signup'>
                    <span class='glyphicon glyphicon-check'></span> 注册
                </a></li>
                -->
            {{end}}
          </ul>
        </li>
      </ul>
    </nav>
    </div>
  </div>

</header>
{{end}}