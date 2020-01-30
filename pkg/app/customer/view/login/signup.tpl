{{ define "login/signup.tpl" }}

{{template "common/head.tpl" .}}
{{template "common/partheader.tpl" .}}

<div class="container">
    <div class="row vertical-offset-50">
    	<div class="col-md-6 col-md-offset-3">
    		<div class="panel panel-default">
			  	<div class="panel-heading text-center">
			    	<h3 class="panel-title"><strong>注册</strong></h3>
			 	</div> 

			  	<div class="panel-body">
			    	<form accept-charset="utf-8" role="form" class="form-horizontal" method="POST" action='/signup'>

                      <div class="form-group">
                        <label class="col-sm-3 control-label">邮箱地址</label>
                        <div class="col-sm-8">
                          <input class="form-control" placeholder="例: admin@example.com" name="email" value="" type="email" required />
                        </div>
                      </div>
                      <div class="form-group">
                        <label class="col-sm-3 control-label">密码</label>
                        <div class="col-sm-8">
			    		            <input class="form-control" placeholder="输入密码" name="pwd" type="password" value="" required
                                    pattern=".{6,}" title="密码长度至少为6个字符" />
                          <input class="form-control" placeholder="确认密码" name="repwd" type="password" required
                                    pattern=".{6,}" title="密码长度至少为6个字符" />
                        </div>
                      </div>
                      <div class="form-group">
                          <label for="inputCaptcha" class="col-sm-3 control-label">验证码</label>
                          <div class="col-sm-4">
                            <input class="form-control" name="captcha" type="text">
                          </div>
                          <img src="{{captchaUrl}}" style="max-width:150px">
                      </div>
                      <div class="form-group text-center">
                          <label for="isCheckedProtocol"><input type="checkbox" id="isCheckedProtocol"/>我同意<a href="javascript:;">注册协议</a></label>
                      </div>
                      <div class="form-group">
                        <div class="col-sm-12">
			    		            <input class="btn btn-lg btn-success btn-block" type="submit" value="注册">
                        </div>
                      </div>
                    </form>
			    </div>

                <div class="panel-footer text-center clearfix">如果您已有账号 <a href='{{urlfor "UsersController.Login"}}'>登录 »</a></div>

			</div>
		</div>
	</div>
</div>
{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}

{{ end }}