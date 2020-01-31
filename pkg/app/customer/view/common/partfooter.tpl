{{define "common/partfooter.tpl"}}
<footer>
  <div class="container">
    <div class="clearfix">
      <div class="footer-logo">
        <a href="/">
          <img src=""><small style='font-size: 65%;'>{{.appname}}</small>
        </a>
      </div>
      <dl class="footer-nav">
        <dt class="nav-title">导航</dt>
        <dd class="nav-item">
          <a href="#">
            <span class="glyphicon glyphicon-home"> 首页</span>
          </a>
        </dd>
        <dd class="nav-item">
          <a href="#" target="_blank">
            <span class='glyphicon glyphicon-briefcase'> 文件</span>
          </a>
        </dd>
        <dd class="nav-item">
          <a href="#" target="_blank">
            <span class='glyphicon glyphicon-upload'> 上传</span>
          </a>
        </dd>
        <dd class="nav-item">
          <a href="#" target="_blank">
            <span class='glyphicon glyphicon-search'> 分享市场</span>
          </a>
        </dd>
      </dl>
      <dl class="footer-nav">
        <dt class="nav-title">相关技术</dt>

        <dd class="nav-item">
          <a href="#">
            <span class='glyphicon glyphicon-info-sign'> Golang</span>
          </a>
        </dd>
        <dd class="nav-item">
          <a href="#">
            <span class='glyphicon glyphicon-info-sign'> gRPC</span>
          </a>
        </dd>
        <dd class="nav-item">
          <a href="#">
            <span class='glyphicon glyphicon-info-sign'> gin</span>
          </a>
        </dd>
        <dd class="nav-item">
          <a href="#">
            <span class='glyphicon glyphicon-info-sign'> Bootstrap</span>
          </a>
        </dd>

      </dl>

      <dl class="footer-nav hidden">
        <dt class="nav-title">合作伙伴</dt>
        <dd class="nav-item">
          <a href="#" target="_blank">
            <span class='glyphicon glyphicon-globe'></span> 
          </a>
        </dd>
      </dl>

      <dl class="footer-nav">
        <dt class="nav-title">联系我们</dt>
        <dd class="nav-item">
          <a href="https://github.com/uxff/flexdrive" target="_blank">
            <span class='glyphicon glyphicon-comment'> view source</span>
          </a>
        </dd>
      </dl>

    </div>

    <div class="footer-copyright text-center">
      Copyright <span class="glyphicon glyphicon-copyright-mark"></span>
      2014-{{datenow "2006"}} <strong>{{.appname}}</strong>
      All rights reserved.
    </div>

  </div>
</footer>
{{end}}