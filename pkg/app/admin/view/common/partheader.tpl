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
        <li><a href='/'>
          <span class="glyphicon glyphicon-home"></span> Fancy Navigator
        </a></li>
        <li>
          <a href="javascript:;" class="dropdown-toggle" data-hover="dropdown">
            <span class="glyphicon glyphicon-picture"></span> Picset <b class="caret"></b>
          </a>
          <ul class="dropdown-menu">
            <li role="presentation" class="dropdown-header">main route</li>
            <li><a href="/picset">Picset</a></li>
            <li role="presentation" class="divider"></li>
            <li role="presentation" class="dropdown-header">selected picsets</li>
            <li><a href="/picset/folderName1/">FolderName1</a></li>
            <li><a href="/picset/55156/">The 55156 site</a></li>
          </ul>
        </li>
      </ul>

      <ul class="nav navbar-nav navbar-right">
        <li class="dropdown">
          <a href="javascript:;" role="button" class="dropdown-toggle" data-hover="dropdown">
            <span class='glyphicon glyphicon-info-sign'></span> Account <b class="caret"></b>
          </a>
          <ul class="dropdown-menu">
            {{if .IsLogin}}
                <li ><a href='{{urlfor "UsersController.Logout"}}'>
                  <span class='glyphicon glyphicon-log-out'></span> Logout
                </a></li>
            {{else}}
                <li ><a href='{{urlfor "UsersController.Login"}}'>
                  <span class='glyphicon glyphicon-log-in'></span> Login
                </a></li>
                <li ><a href='{{urlfor "UsersController.Signup"}}'>
                    <span class='glyphicon glyphicon-check'></span> Sign Up
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