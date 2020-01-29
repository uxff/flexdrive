{{ define "login/changepwd.tpl" }}

{{template "common/head.tpl" .}}
{{template "common/partheader.tpl" .}}

<div class="container">
    <div class="row vertical-offset-75">
    	<div class="col-md-6 col-md-offset-3">
    		<div class="panel panel-default">
			  	<div class="panel-heading text-center">
			    	<h3 class="panel-title"><strong>修改密码</strong></h3>
			 	</div> 

			  	<div class="panel-body">
			    	<form accept-charset="utf-8" role="form" class="form-horizontal" method="POST" action="/changePwd">

                      <div class="form-group">
                        <label for="inputEmail" class="col-sm-3 control-label">输入我的邮箱地址</label>
                        <div class="col-sm-8">
                          <input class="form-control" placeholder="例: admin@example.com" name="email" value="" type="email" required 
                                    id="inputEmail" />
                        </div>
                      </div>
                      <div class="form-group">
                        <label class="col-sm-3 control-label">原密码</label>
                        <div class="col-sm-8">
			    		          <input class="form-control" placeholder="输入秘钥" name="oldpassword" type="password" value="" required
                                    pattern=".{6,}" title="密码长度至少为6个字符"  />
                        </div>
                      </div>
                      <div class="form-group">
                        <label class="col-sm-3 control-label">新密码</label>
                        <div class="col-sm-8">
			    		          <input class="form-control" placeholder="输入秘钥" name="newpassword" type="password" value="" required
                                    pattern=".{6,}" title="密码长度至少为6个字符"  />
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
                        <div class="col-sm-12">
			    		            <input class="btn btn-lg btn-success btn-block" type="submit" value="提交">
                        </div>
                      </div>
                    </form>
			    </div>

          <div class="panel-footer text-center clearfix" style="display:none">没有账户 <a href="/signup">注册 »</a></div>

			</div>
		</div>
	</div>
</div>
{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}

{{ end }}