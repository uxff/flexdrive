{{define "common/partheader.tpl"}}

<header id="topbar" class="navbar navbar-inverse bs-docs-nav" role="banner">
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
        <strong>分布式云盘系统</strong>
      </a>
    </div>

    <nav class="collapse navbar-collapse bs-navbar-collapse" role="navigation" >
      <!-- 左侧主要菜单 -->
      <ul class="nav navbar-nav">
        <li>
          <a href='/'>
            <span class="glyphicon glyphicon-home"></span> 首页
          </a>
        </li>
        <li>
          <a href="/my/file/list" class="dropdown-toggle" data-hover="dropdown" >
            <span class="glyphicon glyphicon-briefcase"></span> 文件 
          </a>
        </li>
        <li class="hidden">
          <a href="/upload" class="dropdown-toggle" data-hover="dropdown">
            <span class="glyphicon glyphicon-cloud-upload"></span> 上传 
          </a>
        </li>
        <li>
          <a href="/share/search" class="dropdown-toggle" data-hover="dropdown" >
            <span class="glyphicon glyphicon-search"></span> 分享市场 
          </a>
        </li>
        <li>
          <a href="javascript:;" class="dropdown-toggle" data-hover="dropdown">
            <span class="glyphicon glyphicon-cog"></span> 我的 <b class="caret"></b>
          </a>
          <ul class="dropdown-menu">
            <li><a href="/my/profile">总览</a></li>
            <li><a href="/my/share">我的分享</a></li>
            <li><a href="/my/downloadtask">我的离线下载</a></li>
            <li><a href="/my/order/list">我的订单</a></li>
          </ul>
        </li>
      </ul>

      <!-- 右侧菜单 -->
      <ul class="nav navbar-nav navbar-right">
        <li class="dropdown">
          <a href="javascript:;" role="button" class="dropdown-toggle" data-hover="dropdown">
            <span class='glyphicon glyphicon-user'></span> 
            {{if .IsLogin}}
              {{.LoginInfo.UserEnt.Email}} 
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
            {{end}}
              <li ><a href='/signup'>
                  <span class='glyphicon glyphicon-check'></span> 注册
              </a></li>
          </ul>
        </li>
      </ul>
    </nav>
    </div>
  </div>

</header>
{{end}}