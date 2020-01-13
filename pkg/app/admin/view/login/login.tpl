{{append . "HeadStyles" "/static/css/custom.css"}}
{{append . "HeadScripts" "/static/js/custom.js"}}

<div class="container">
    <div class="row vertical-offset-75">
    	<div class="col-md-6 col-md-offset-3">
    		<div class="panel panel-default">
			  	<div class="panel-heading text-center">
			    	<h3 class="panel-title"><strong>登录</strong></h3>
			 	</div> 

			  	<div class="panel-body">
			    	<form accept-charset="utf-8" role="form" class="form-horizontal" method="POST" action='{{urlfor "UsersController.Login"}}'>
                      {{ .xsrfdata }}

                      {{template "alert.tpl" .}}

                      <div class="form-group">
                        <label for="inputEmail" class="col-sm-3 control-label">邮箱地址</label>
                        <div class="col-sm-8">
                          <input class="form-control" placeholder="例: admin@example.com" name="Email" value="{{index .Params "Email"}}" type="email" required 
                                    id="inputEmail" />
                        </div>
                      </div>
                      <div class="form-group">
                        <label for="inputPassword" class="col-sm-3 control-label">秘钥</label>
                        <div class="col-sm-8">
			    		  <input class="form-control" placeholder="输入秘钥" name="Password" type="password" value="" required
                                    pattern=".{6,}" title="密码长度至少为6个字符" id="inputPassword"  />
                        </div>
                      </div>
                      <div class="form-group">
                          <label for="inputCaptcha" class="col-sm-3 control-label">验证码</label>
                          <div class="col-sm-4">
                            <input class="form-control" name="captcha" type="text">
                          </div>
                        {..{ create_captcha}.. }
                      </div>
                      <div class="form-group text-center">
                        <div class="col-sm-12">
			    		  <input class="btn btn-lg btn-success btn-block" type="submit" value="登录">
                            <a href="{urlfor "UsersController.PasswordReset"}">
                                忘记秘钥，请点击此处 »
                            </a>

                        </div>
                      </div>
                    </form>
			    </div>

                <div class="panel-footer text-center clearfix">没有账户 <a href='{urlfor "UsersController.Signup"}'>注册 »</a></div>

			</div>
		</div>
	</div>
</div>
