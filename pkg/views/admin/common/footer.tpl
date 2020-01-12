<footer>
  <div class="container">
    <div class="clearfix">
      <div class="footer-logo">
        <a href="/">
          <img src=""><small style='font-size: 65%;'>{{.appname}}</small>
        </a>
      </div>
      <dl class="footer-nav">
        <dt class="nav-title">GALLERY</dt>
        <dd class="nav-item">
          <a href="#">
            <span class="glyphicon glyphicon-credit-card"> Donate</span>
          </a>
        </dd>
        <dd class="nav-item">
          <a href="#" target="_blank">
            <span class='glyphicon glyphicon-bullhorn'> Present</span>
          </a>
        </dd>
      </dl>
      <dl class="footer-nav">
        <dt class="nav-title">ABOUT</dt>

        <dd class="nav-item">
          <a href="#">
            <span class='glyphicon glyphicon-info-sign'> </span>
          </a>
        </dd>

      </dl>

      <dl class="footer-nav hidden">
        <dt class="nav-title">SOCIAL</dt>
        <dd class="nav-item">
          <a href="#" target="_blank">
            <span class='glyphicon glyphicon-globe'></span> 
          </a>
        </dd>
      </dl>

      <dl class="footer-nav">
        <dt class="nav-title">CONTACT</dt>
        <dd class="nav-item">
          <a href="#">
            <span class='glyphicon glyphicon-comment'></span>
          </a>
        </dd>
      </dl>

    </div>

    <div class="footer-copyright text-center">
      友情链接:
      {{range $k, $link := .friendlyLinks}}
          <a href="{{$link.Url}}" target="_blank">{{$link.Name}}</a> &nbsp;
      {{end}}
    </div>
    <div class="footer-copyright text-center">
      Copyright <span class="glyphicon glyphicon-copyright-mark"></span>
      2014-{{datenow "2006"}} <strong>{{.appname}}</strong>
      All rights reserved.
    </div>

  </div>
</footer>
