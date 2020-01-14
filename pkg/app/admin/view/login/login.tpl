{{ define "login/login.tpl" }}

{{template "common/head.tpl" .}}
{{template "common/partheader.tpl" .}}

<div class="container">
    <div class="row vertical-offset-75">
    	<div class="col-md-6 col-md-offset-3">
    		<div class="panel panel-default">
			  	<div class="panel-heading text-center">
			    	<h3 class="panel-title"><strong>登录</strong></h3>
			 	</div> 

			  	<div class="panel-body">


                      <div class="form-group">
                        <label for="inputEmail" class="col-sm-3 control-label">邮箱地址</label>
                        <div class="col-sm-8">
                        </div>
                      </div>
                      <div class="form-group">
                        <label for="inputPassword" class="col-sm-3 control-label">秘钥</label>
                        <div class="col-sm-8">
                        </div>
                      </div>
                      <div class="form-group">
                          <label for="inputCaptcha" class="col-sm-3 control-label">验证码</label>
                          <div class="col-sm-4">
                            <input class="form-control" name="captcha" type="text">
                          </div>
                      </div>
                      <div class="form-group text-center">
                        <div class="col-sm-12">
			    		  <input class="btn btn-lg btn-success btn-block" type="submit" value="登录">

                        </div>
                      </div>
			    </div>


			</div>
		</div>
	</div>
</div>

{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}

{{ end }}