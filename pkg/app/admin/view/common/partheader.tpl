{{define "common/partheader.tpl"}}
<header id="topbar" class="navbar navbar-default navbar-fixed-top bs-docs-nav" role="banner">
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
      <ul class="nav navbar-nav">
        <li>
          <a href='/'>
            <span class="glyphicon glyphicon-home"></span> Fancy Navigator
          </a>
        </li>
        <li>
          <a href="javascript:;" class="dropdown-toggle" data-hover="dropdown">
            <span class="glyphicon glyphicon-picture"></span> 节点管理 <b class="caret"></b>
          </a>
          <ul class="dropdown-menu">
            <li role="presentation" class="dropdown-header">main route</li>
            <li><a href="/picset">节点列表</a></li>
            <li role="presentation" class="divider"></li>
            <li role="presentation" class="dropdown-header">selected picsets</li>
            <li><a href="/picset/folderName1/">文件管理</a></li>
            <li><a href="/picset/55156/"></a></li>
          </ul>
        </li>
        <li>
          <a href="javascript:;" class="dropdown-toggle" data-hover="dropdown">
            <span class="glyphicon glyphicon-picture"></span> 会员管理 <b class="caret"></b>
          </a>
          <ul class="dropdown-menu">
            <li><a href="/picset">会员列表</a></li>
            <li><a href="/picset/folderName1/">等级管理</a></li>
            <li><a href="/picset/55156/">订单管理</a></li>
            <li><a href="/picset/55156/">分享管理</a></li>
          </ul>
        </li>
        <li>
          <a href="javascript:;" class="dropdown-toggle" data-hover="dropdown">
            <span class="glyphicon glyphicon-picture"></span> 系统管理 <b class="caret"></b>
          </a>
          <ul class="dropdown-menu">
            <li role="presentation" class="dropdown-header">管理员账号</li>
            <li><a href="/picset">Picset</a></li>
            <li role="presentation" class="divider"></li>
            <li role="presentation" class="dropdown-header">selected picsets</li>
            <li><a href="/picset/folderName1/">FolderName1</a></li>
            <li><a href="/picset/55156/">The 55156 site</a></li>
          </ul>
        </li>
      </ul>

      <!-- 右侧菜单 -->
      <ul class="nav navbar-nav navbar-right">
        <li class="dropdown">
          <a href="javascript:;" role="button" class="dropdown-toggle" data-hover="dropdown">
            <span class='glyphicon glyphicon-info-sign'></span> 账户 <b class="caret"></b>
          </a>
          <ul class="dropdown-menu">
            {{if .IsLogin}}
                <li ><a href='{{urlfor "UsersController.Logout"}}'>
                  <span class='glyphicon glyphicon-log-out'></span> 退出登录
                </a></li>
            {{else}}
                <li ><a href='{{urlfor "UsersController.Login"}}'>
                  <span class='glyphicon glyphicon-log-in'></span> 登录
                </a></li>
                <li ><a href='{{urlfor "UsersController.Signup"}}'>
                    <span class='glyphicon glyphicon-check'></span> 注册
                </a></li>
            {{end}}
          </ul>
        </li>
      </ul>
    </nav>
    </div>
  </div>

</header>
{{end}}